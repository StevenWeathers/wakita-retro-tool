<script>
    import Sockette from 'sockette'
    import { onMount, onDestroy } from 'svelte'
    import dragula from 'dragula'

    import PageLayout from '../components/PageLayout.svelte'
    import UserCard from '../components/UserCard.svelte'
    import InviteUser from '../components/InviteUser.svelte'
    import UsersIcon from '../components/icons/UsersIcon.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import DownCarrotIcon from '../components/icons/DownCarrotIcon.svelte'
    import ChevronRight from '../components/icons/ChevronRight.svelte'
    import DeleteRetrospective from '../components/DeleteRetrospective.svelte'
    import SolidButton from '../components/SolidButton.svelte'
    import CheckCircle from '../components/icons/CheckCircle.svelte'
    import CheckboxIcon from '../components/icons/CheckboxIcon.svelte'
    import CrossCircle from '../components/icons/CrossCircle.svelte'
    import RetroItemForm from '../components/RetroItemForm.svelte'
    import ArrowUp from '../components/icons/ArrowUp.svelte'
    import { appRoutes, PathPrefix } from '../config'
    import { user } from '../stores.js'

    export let retrospectiveId
    export let notifications
    export let router
    export let eventTag

    const { AllowRegistration } = appConfig
    const loginOrRegister = AllowRegistration
        ? appRoutes.register
        : appRoutes.login

    const hostname = window.location.origin
    const socketExtension = window.location.protocol === 'https:' ? 'wss' : 'ws'

    let socketError = false
    let socketReconnecting = false
    let retrospective = {
        ownerId: '',
        phase: 1,
        users: [],
        workedItems: [],
        improveItems: [],
        questionItems: [],
        actionItems: [],
    }
    let showUsers = false
    let showDeleteRetrospective = false
    let actionItem = ''
    let showExport = false

    const onSocketMessage = function(evt) {
        const parsedEvent = JSON.parse(evt.data)

        switch (parsedEvent.type) {
            case 'init':
                retrospective = JSON.parse(parsedEvent.value)
                retrospective.workedItems = nestItems(retrospective.workedItems)
                retrospective.improveItems = nestItems(
                    retrospective.improveItems,
                )
                retrospective.questionItems = nestItems(
                    retrospective.questionItems,
                )
                eventTag('join', 'retrospective', '')
                break
            case 'user_joined': {
                retrospective.users = JSON.parse(parsedEvent.value)
                const joinedUser = retrospective.users.find(
                    w => w.id === parsedEvent.userId,
                )
                notifications.success(`${joinedUser.name} joined.`)
                break
            }
            case 'user_retreated': {
                const leftUser = retrospective.users.find(
                    w => w.id === parsedEvent.userId,
                )
                retrospective.users = JSON.parse(parsedEvent.value)

                notifications.danger(`${leftUser.name} retreated.`)
                break
            }
            case 'retrospective_updated':
                retrospective = JSON.parse(parsedEvent.value)
                retrospective.workedItems = nestItems(retrospective.workedItems)
                retrospective.improveItems = nestItems(
                    retrospective.improveItems,
                )
                retrospective.questionItems = nestItems(
                    retrospective.questionItems,
                )
                break
            case 'item_worked_updated': {
                const parsedValue = JSON.parse(parsedEvent.value)
                retrospective.workedItems = nestItems(parsedValue)
                break
            }
            case 'item_improve_updated': {
                const parsedValue = JSON.parse(parsedEvent.value)
                retrospective.improveItems = nestItems(parsedValue)
                break
            }
            case 'item_question_updated': {
                const parsedValue = JSON.parse(parsedEvent.value)
                retrospective.questionItems = nestItems(parsedValue)
                break
            }
            case 'action_updated':
                retrospective.actionItems = JSON.parse(parsedEvent.value)
                break
            case 'retrospective_conceded':
                // retrospective over, goodbye.
                notifications.warning('Retrospective deleted')
                router.route(appRoutes.retrospectives)
                break
            default:
                break
        }
    }

    const ws = new Sockette(
        `${socketExtension}://${window.location.host}${PathPrefix}/api/arena/${retrospectiveId}`,
        {
            timeout: 2e3,
            maxAttempts: 15,
            onmessage: onSocketMessage,
            onerror: () => {
                socketError = true
                eventTag('socket_error', 'retrospective', '')
            },
            onclose: e => {
                if (e.code === 4004) {
                    eventTag('not_found', 'retrospective', '', () => {
                        router.route(appRoutes.retrospectives)
                    })
                } else if (e.code === 4001) {
                    eventTag('socket_unauthorized', 'retrospective', '', () => {
                        user.delete()
                        router.route(`${appRoutes.login}/${retrospectiveId}`)
                    })
                } else if (e.code === 4003) {
                    eventTag('socket_duplicate', 'retrospective', '', () => {
                        notifications.danger(
                            `Duplicate retrospective session exists for your ID`,
                        )
                        router.route(`${appRoutes.retrospectives}`)
                    })
                } else if (e.code === 4002) {
                    eventTag(
                        'retrospective_user_abandoned',
                        'retrospective',
                        '',
                        () => {
                            router.route(appRoutes.retrospectives)
                        },
                    )
                } else {
                    socketReconnecting = true
                    eventTag('socket_close', 'retrospective', '')
                }
            },
            onopen: () => {
                socketError = false
                socketReconnecting = false
                eventTag('socket_open', 'retrospective', '')
            },
            onmaximum: () => {
                socketReconnecting = false
                eventTag(
                    'socket_error',
                    'retrospective',
                    'Socket Reconnect Max Reached',
                )
            },
        },
    )

    onDestroy(() => {
        eventTag('leave', 'retrospective', '', () => {
            ws.close()
        })
    })

    $: isOwner = retrospective.ownerId === $user.id

    const sendSocketEvent = (type, value) => {
        ws.send(
            JSON.stringify({
                type,
                value,
            }),
        )
    }

    function concedeRetrospective() {
        eventTag('concede_retrospective', 'retrospective', '', () => {
            sendSocketEvent('concede_retrospective', '')
        })
    }

    function abandonRetrospective() {
        eventTag('abandon_retrospective', 'retrospective', '', () => {
            sendSocketEvent('abandon_retrospective', '')
        })
    }

    function toggleUsersPanel() {
        showUsers = !showUsers
        eventTag('show_users', 'retrospective', `show: ${showUsers}`)
    }

    const toggleDeleteRetrospective = () => {
        showDeleteRetrospective = !showDeleteRetrospective
    }

    const toggleExport = () => {
        showExport = !showExport
    }

    const handleItemAdd = (type, content) => {
        sendSocketEvent(
            `create_item_${type}`,
            JSON.stringify({
                content,
                phase: retrospective.phase,
            }),
        )
    }

    const unnestItem = (type, id) => () => {
        sendSocketEvent(
            `unnest_item_${type}`,
            JSON.stringify({
                id,
                parentId: '',
            }),
        )
    }

    const handleItemDelete = (type, id) => () => {
        sendSocketEvent(
            `delete_item_${type}`,
            JSON.stringify({
                id,
                phase: retrospective.phase,
            }),
        )
    }

    const handleActionItem = evt => {
        evt.preventDefault()

        sendSocketEvent(
            'create_action',
            JSON.stringify({
                content: actionItem,
            }),
        )
        actionItem = ''
    }

    const handleActionUpdate = (id, completed) => evt => {
        sendSocketEvent(
            'update_action',
            JSON.stringify({
                id,
                completed: !completed,
            }),
        )
    }

    const handleActionDelete = id => () => {
        sendSocketEvent(
            'delete_action',
            JSON.stringify({
                id,
            }),
        )
    }

    const advancePhase = () => {
        sendSocketEvent(
            'advance_phase',
            JSON.stringify({
                phase: retrospective.phase + 1,
            }),
        )
    }

    const voteItem = (type, id) => {
        sendSocketEvent(
            `vote_item_${type}`,
            JSON.stringify({
                id,
            }),
        )
    }

    const nestItems = items => {
        const parentMap = {}
        const nestedItems = items.reduce((prev, item) => {
            if (item.parentId !== '') {
                parentMap[item.parentId] = parentMap[item.parentId] || {
                    items: [],
                    voteCount: 0,
                }
                parentMap[item.parentId].voteCount =
                    parentMap[item.parentId].voteCount + item.votes.length
                parentMap[item.parentId].items.push(item)
                return prev
            }
            prev.push(item)

            return prev
        }, [])

        nestedItems.forEach(item => {
            item.voteCount = item.votes.length
            if (typeof parentMap[item.id] !== 'undefined') {
                item.items = parentMap[item.id].items
                item.voteCount = item.voteCount + parentMap[item.id].voteCount
            } else {
                item.items = []
            }
        })

        return nestedItems
    }

    let drake = dragula({
        isContainer: function(el) {
            return el.classList.contains('item-list-item')
        },
        moves: function(el, source, handle, sibling) {
            return true
        },
        accepts: function(el, target, source) {
            return (
                source.dataset.itemtype === target.dataset.itemtype &&
                target !== source
            )
        },
        invalid: function(el) {
            const isDisabled = el.dataset.dragdisabled || 'false'
            return (
                isDisabled === 'true' || retrospective.phase === 1 || !isOwner
            )
        },
    })

    drake.on('drop', function(el, target, source) {
        const itemType = source.dataset.itemtype
        const itemId = source.dataset.itemid
        const parentItemId = target.dataset.itemid

        el.remove()
        sendSocketEvent(
            `nest_item_${itemType}`,
            JSON.stringify({
                id: itemId,
                parentId: parentItemId,
            }),
        )
    })

    onMount(() => {
        if (!$user.id) {
            router.route(`${loginOrRegister}/${retrospectiveId}`)
        }
    })
</script>

<style>
    /** Manually including Dragula styles, should automate this later */
    :global(.gu-mirror) {
        position: fixed !important;
        margin: 0 !important;
        z-index: 9999 !important;
        opacity: 0.8;
        -ms-filter: 'progid:DXImageTransform.Microsoft.Alpha(Opacity=80)';
        filter: alpha(opacity=80);
    }

    :global(.gu-hide) {
        display: none !important;
    }

    :global(.gu-unselectable) {
        -webkit-user-select: none !important;
        -moz-user-select: none !important;
        -ms-user-select: none !important;
        user-select: none !important;
    }

    :global(.gu-transit) {
        opacity: 0.2;
        -ms-filter: 'progid:DXImageTransform.Microsoft.Alpha(Opacity=20)';
        filter: alpha(opacity=20);
    }
    :global(input:checked ~ div) {
        @apply border-green-500;
    }
    :global(input:checked ~ div svg) {
        @apply block;
    }
</style>

<svelte:head>
    <title>Retrospective {retrospective.name} | Wakita</title>
</svelte:head>

{#if retrospective.name && !socketReconnecting && !socketError}
    <div
        class="px-6 py-2 bg-gray-100 border-b border-t border-gray-400 flex
        flex-wrap">
        <div class="w-1/4">
            <h1 class="text-3xl font-bold leading-tight">
                {retrospective.name}
            </h1>
        </div>
        <div class="w-3/4 text-right">
            <div>
                {#if retrospective.phase === 4}
                    <SolidButton color="green" onClick="{toggleExport}">
                        {#if showExport}
                            Back
                        {:else}
                            Export
                            <ArrowUp
                                class="inline-block ml-1"
                                width="12"
                                height="12" />
                        {/if}
                    </SolidButton>
                {/if}
                {#if isOwner}
                    {#if retrospective.phase !== 4}
                        <SolidButton color="blue" onClick="{advancePhase}">
                            {#if retrospective.phase === 1}
                                Group &amp; Vote comments
                            {:else if retrospective.phase === 2}
                                Discuss and add action items
                            {:else if retrospective.phase === 3}
                                Finish retro
                            {/if}
                        </SolidButton>
                    {/if}

                    <HollowButton
                        color="red"
                        onClick="{toggleDeleteRetrospective}"
                        class="mr-2">
                        Delete Retrospective
                    </HollowButton>
                {:else}
                    <HollowButton color="red" onClick="{abandonRetrospective}">
                        Leave Retrospective
                    </HollowButton>
                {/if}
                <div class="inline-block relative">
                    <HollowButton
                        color="gray"
                        class="transition ease-in-out duration-150"
                        onClick="{toggleUsersPanel}">
                        <UsersIcon class="mr-1" height="18" width="18" />
                        Users
                        <DownCarrotIcon class="ml-1" />
                    </HollowButton>
                    {#if showUsers}
                        <div
                            class="origin-top-right absolute right-0 mt-1 w-64
                            rounded-md shadow-lg text-left">
                            <div class="rounded-md bg-white shadow-xs">
                                {#each retrospective.users as usr, index (usr.id)}
                                    {#if usr.active}
                                        <UserCard
                                            user="{usr}"
                                            showBorder="{index != retrospective.users.length - 1}" />
                                    {/if}
                                {/each}

                                <div class="p-2">
                                    <InviteUser
                                        {hostname}
                                        retrospectiveId="{retrospective.id}" />
                                </div>
                            </div>
                        </div>
                    {/if}
                </div>
            </div>
        </div>
    </div>
    <div class="px-6 py-2 bg-gray-100 border-b border-gray-400 flex flex-wrap">
        <div class="w-1/2">
            <div class="flex items-center text-gray-500">
                <div
                    class="flex-initial px-1 {retrospective.phase === 1 && 'border-b-2 border-blue-500 text-gray-800'}">
                    Brainstorm
                </div>
                <div class="flex-initial px-1">
                    <ChevronRight />
                </div>
                <div
                    class="flex-initial px-1 {retrospective.phase === 2 && 'border-b-2 border-blue-500 text-gray-800'}">
                    Group &amp; Vote
                </div>
                <div class="flex-initial px-1">
                    <ChevronRight />
                </div>
                <div
                    class="flex-initial px-1 {retrospective.phase === 3 && 'border-b-2 border-blue-500 text-gray-800'}">
                    Add action items
                </div>
                <div class="flex-initial px-1">
                    <ChevronRight />
                </div>
                <div
                    class="flex-initial px-1 {retrospective.phase === 4 && 'border-b-2 border-blue-500 text-gray-800'}">
                    Done
                </div>
            </div>
        </div>
        <div class="w-1/2 text-right text-gray-600">
            {#if retrospective.phase === 1}
                Add your comments below, you won't be able to see your peers
                until next step
            {:else if retrospective.phase === 2}
                Drag and drop comments to group them together and vote for the
                ones you'd like to discuss about
            {:else if retrospective.phase === 3}
                Add action items, you can no longer group or vote comments
            {/if}
        </div>

    </div>
    <div class="flex flex-grow p-4">
        {#if showExport}
            <div class="px-4">
                <div class="mb-4">
                    <h2 class="text-2xl font-bold">Works</h2>
                    <ul class="pl-12 list-disc">
                        {#each retrospective.workedItems as item (item.id)}
                            <li>
                                {item.content} ({item.voteCount})
                                {#if item.items.length}
                                    <ul class="pl-8 list-disc">
                                        {#each item.items as child (child.id)}
                                            <li>{child.content}</li>
                                        {/each}
                                    </ul>
                                {/if}
                            </li>
                        {/each}
                    </ul>
                </div>
                <div class="mb-4">
                    <h2 class="text-2xl font-bold">Needs Improvement</h2>
                    <ul class="pl-12 list-disc">
                        {#each retrospective.improveItems as item (item.id)}
                            <li>
                                {item.content} ({item.voteCount})
                                {#if item.items.length}
                                    <ul class="pl-8 list-disc">
                                        {#each item.items as child (child.id)}
                                            <li>{child.content}</li>
                                        {/each}
                                    </ul>
                                {/if}
                            </li>
                        {/each}
                    </ul>
                </div>
                <div class="mb-4">
                    <h2 class="text-2xl font-bold">Questions</h2>
                    <ul class="pl-12 list-disc">
                        {#each retrospective.questionItems as item (item.id)}
                            <li>
                                {item.content} ({item.voteCount})
                                {#if item.items.length}
                                    <ul class="pl-8 list-disc">
                                        {#each item.items as child (child.id)}
                                            <li>{child.content}</li>
                                        {/each}
                                    </ul>
                                {/if}
                            </li>
                        {/each}
                    </ul>
                </div>
                <div class="mb-4">
                    <h2 class="text-2xl font-bold">Action Items</h2>
                    <ul class="pl-12 list-disc">
                        {#each retrospective.actionItems as item (item.id)}
                            <li>{item.content}</li>
                        {/each}
                    </ul>
                </div>
            </div>
        {:else}
            <RetroItemForm
                handleSubmit="{handleItemAdd}"
                handleDelete="{handleItemDelete}"
                handleVote="{voteItem}"
                handleUnnest="{unnestItem}"
                itemType="worked"
                newItemPlaceholder="What worked well..."
                phase="{retrospective.phase}"
                {isOwner}
                items="{retrospective.workedItems}" />
            <RetroItemForm
                handleSubmit="{handleItemAdd}"
                handleDelete="{handleItemDelete}"
                handleVote="{voteItem}"
                handleUnnest="{unnestItem}"
                itemType="improve"
                newItemPlaceholder="What needs improvement..."
                phase="{retrospective.phase}"
                {isOwner}
                items="{retrospective.improveItems}" />
            <RetroItemForm
                handleSubmit="{handleItemAdd}"
                handleDelete="{handleItemDelete}"
                handleVote="{voteItem}"
                handleUnnest="{unnestItem}"
                itemType="question"
                newItemPlaceholder="I want to ask..."
                phase="{retrospective.phase}"
                {isOwner}
                items="{retrospective.questionItems}" />
            <div class="w-1/4 mx-2 p-4 bg-white shadow">
                <div class="flex items-center mb-2">
                    <div class="flex-shrink pr-1">
                        <CheckCircle
                            class="text-gray-400"
                            height="24"
                            width="24" />
                    </div>
                    <div class="flex-grow">
                        <form on:submit="{handleActionItem}">
                            <input
                                bind:value="{actionItem}"
                                placeholder="Action item..."
                                class="border-gray-300 border-2 appearance-none
                                rounded w-full py-2 px-3 text-gray-700
                                leading-tight focus:outline-none focus:bg-white
                                focus:border-orange-500"
                                id="actionItem"
                                name="actionItem"
                                type="text"
                                required
                                disabled="{retrospective.phase !== 3 || !isOwner}" />
                            <button type="submit" class="hidden"></button>
                        </form>
                    </div>
                </div>
                {#each retrospective.actionItems as item, i}
                    <div class="py-1 my-1">
                        <div class="flex content-center">
                            <div class="flex-shrink">
                                {#if isOwner}
                                    <button
                                        on:click="{handleActionDelete(item.id)}"
                                        class="pr-2 pt-1 text-gray-500
                                        hover:text-red-500">
                                        <CrossCircle height="18" width="18" />
                                    </button>
                                {/if}
                            </div>
                            <div class="flex-grow">{item.content}</div>
                            <div class="flex-shrink">
                                <input
                                    type="checkbox"
                                    id="{i}Completed"
                                    checked="{item.completed}"
                                    class="opacity-0 absolute h-6 w-6"
                                    on:change="{handleActionUpdate(item.id, item.completed)}" />
                                <div
                                    class="bg-white border-2 rounded-md
                                    border-gray-400 w-6 h-6 flex flex-shrink-0
                                    justify-center items-center mr-2
                                    focus-within:border-blue-500">
                                    <CheckboxIcon />
                                </div>
                                <label
                                    for="{i}Completed"
                                    class="select-none"></label>
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        {/if}
    </div>
{:else}
    <PageLayout>
        <div class="flex items-center">
            <div class="flex-1 text-center">
                {#if socketReconnecting}
                    <h1
                        class="text-5xl text-orange-500 leading-tight font-bold">
                        Ooops, reloading Retrospective...
                    </h1>
                {:else if socketError}
                    <h1 class="text-5xl text-red-500 leading-tight font-bold">
                        Error joining retrospective, refresh and try again.
                    </h1>
                {:else}
                    <h1 class="text-5xl text-green-500 leading-tight font-bold">
                        Loading Retrospective...
                    </h1>
                {/if}
            </div>
        </div>
    </PageLayout>
{/if}

{#if showDeleteRetrospective}
    <DeleteRetrospective
        toggleDelete="{toggleDeleteRetrospective}"
        handleDelete="{concedeRetrospective}" />
{/if}
