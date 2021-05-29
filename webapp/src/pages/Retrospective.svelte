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
    import QuestionCircle from '../components/icons/QuestionCircle.svelte'
    import CheckCircle from '../components/icons/CheckCircle.svelte'
    import SmileCircle from '../components/icons/SmileCircle.svelte'
    import FrownCircle from '../components/icons/FrownCircle.svelte'
    import CheckboxIcon from '../components/icons/CheckboxIcon.svelte'
    import { appRoutes, PathPrefix } from '../config'
    import { user } from '../stores.js'
import CrossCircle from '../components/icons/CrossCircle.svelte'
import ArrowLeft from '../components/icons/ArrowLeft.svelte'

    export let retrospectiveId
    export let notifications
    export let router
    export let eventTag

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
        actionItems: []
    }
    let showUsers = false
    let showDeleteRetrospective = false

    let workedWell = ''
    let needsImprovement = ''
    let question = ''
    let actionItem = ''

    const onSocketMessage = function(evt) {
        const parsedEvent = JSON.parse(evt.data)

        switch (parsedEvent.type) {
            case 'init':
                retrospective = JSON.parse(parsedEvent.value)
                if (retrospective.phase !== 1) {
                    retrospective.workedItems = nestItems(retrospective.workedItems)
                    retrospective.improveItems = nestItems(retrospective.improveItems)
                    retrospective.questionItems = nestItems(retrospective.questionItems)
                }
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
                break
            case 'item_worked_updated': {
                const parsedValue = JSON.parse(parsedEvent.value)
                retrospective.workedItems = retrospective.phase !== 1 ? nestItems(parsedValue) : parsedValue
                break;
            }
            case 'item_improve_updated': {
                const parsedValue = JSON.parse(parsedEvent.value)
                retrospective.improveItems = retrospective.phase !== 1 ? nestItems(parsedValue) : parsedValue
                break;
            }
            case 'item_question_updated': {
                const parsedValue = JSON.parse(parsedEvent.value)
                retrospective.questionItems = retrospective.phase !== 1 ? nestItems(parsedValue) : parsedValue
                break;
            }
            case 'action_updated':
                retrospective.actionItems = JSON.parse(parsedEvent.value)
                break;
            case 'phase_updated':
                retrospective.phase =  JSON.parse(parsedEvent.value)
                break;
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

    const handleWorkedWell = (evt) => {
        evt.preventDefault()

        sendSocketEvent('create_item_worked', JSON.stringify({
            content: workedWell,
            phase: retrospective.phase
        }))
        workedWell = ''
    }

    const handleNeedsImprovement = (evt) => {
        evt.preventDefault()

        sendSocketEvent('create_item_improve', JSON.stringify({
            content: needsImprovement,
            phase: retrospective.phase
        }))
        needsImprovement = ''
    }

    const handleQuestion = (evt) => {
        evt.preventDefault()

        sendSocketEvent('create_item_question', JSON.stringify({
            content: question,
            phase: retrospective.phase
        }))
        question = ''
    }

    const unnestItem = (type, id) => () => {
        sendSocketEvent(`unnest_item_${type}`, JSON.stringify({
            id,
            parentId: ''
        }))
    }

    const handleItemDelete = (type, id) => () => {
        sendSocketEvent(`delete_item_${type}`, JSON.stringify({
            id,
            phase: retrospective.phase
        }))
    }

    const handleActionItem = (evt) => {
        evt.preventDefault()

        sendSocketEvent('create_action', JSON.stringify({
            content: actionItem
        }))
        actionItem = ''
    }

    const handleActionUpdate = (id, completed) => (evt) => {
        console.log(id)
        console.log(completed)

        sendSocketEvent('update_action', JSON.stringify({
            id,
            completed: !completed
        }))
    }

    const handleActionDelete = (id) => () => {
        sendSocketEvent('delete_action', JSON.stringify({
            id
        }))
    }

    const advancePhase = () => {
        sendSocketEvent('advance_phase', JSON.stringify({
            phase: retrospective.phase + 1
        }))
    }

    const nestItems = (items) => {
        const parentMap = {}
        const nestedItems = items.reduce((prev, item) => {
            if (item.parentId !== "") {
                parentMap[item.parentId] = parentMap[item.parentId] || []
                parentMap[item.parentId].push(item)
                return prev
            }
            prev.push(item)

            return prev
        }, [])

        nestedItems.forEach(item => {
            if (typeof parentMap[item.id] !== 'undefined') {
                item.items = parentMap[item.id]
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
        moves: function (el, source, handle, sibling) {
            return true
        },
        accepts: function (el, target, source) {
            return source.dataset.itemtype === target.dataset.itemtype && target !== source
        },
        invalid: function (el) {
            const isDisabled = el.dataset.dragdisabled || "false"
            return isDisabled === "true" || retrospective.phase === 1 || !isOwner
        }
    })

    drake.on('drop', function(el, target, source, sibling) {
        const itemType = source.dataset.itemtype
        const itemTypeKey = `${itemType}Items`
        const itemId = source.dataset.itemid
        const parentItemId = target.dataset.itemid
        const childItem = retrospective[itemTypeKey].find(item => item.id === itemId)

        el.remove()
        sendSocketEvent(`nest_item_${itemType}`, JSON.stringify({
            id: itemId,
            parentId: parentItemId
        }))
    })

    onMount(() => {
        if (!$user.id) {
            router.route(`${appRoutes.login}/${retrospectiveId}`)
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
    <div class="px-6 py-2 bg-gray-100 border-b border-t border-gray-400 flex flex-wrap">
        <div class="w-1/4">
            <h1 class="text-3xl font-bold leading-tight">
                {retrospective.name}
            </h1>
        </div>
        <div class="w-3/4 text-right">
            <div>
                {#if isOwner}
                    {#if retrospective.phase !== 4}
                        <SolidButton color="blue" onClick={advancePhase}>
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
                        additionalClasses="mr-2">
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
                        additionalClasses="transition ease-in-out duration-150"
                        onClick="{toggleUsersPanel}">
                        <UsersIcon
                            additionalClasses="mr-1"
                            height="18"
                            width="18" />
                        Users
                        <DownCarrotIcon additionalClasses="ml-1" />
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
        <div class="w-1/3">
            Brainstorm <ChevronRight class="inline-block" /> Group &amp; Vote <ChevronRight class="inline-block" /> Add action items <ChevronRight class="inline-block" /> Done
        </div>
        <div class="w-2/3 text-right text-gray-600">
            {#if retrospective.phase === 1}
                Add your comments below, you won't be able to see your peers until next step
            {:else if retrospective.phase === 2}
                Drag and drop comments to group them together and vote for the ones you'd like to discuss about
            {:else if retrospective.phase === 3}
                Add action items, you can no longer group or vote comments
            {/if}
        </div>

    </div>
    <div class="flex p-4 min-h-screen">
        <div class="w-1/4 mx-2 p-4 bg-white shadow">
            <div class="flex items-center mb-2">
                <div class="flex-shrink pr-1">
                    <SmileCircle class="text-gray-400" height="24" width="24" />
                </div>
                <div class="flex-grow">
                    <form on:submit={handleWorkedWell} class="flex">
                        <input
                            bind:value="{workedWell}"
                            placeholder="What worked well..."
                            class="border-gray-300 border-2
                            appearance-none rounded py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-orange-500 w-full"
                            id="workedWell"
                            name="workedWell"
                            type="text"
                            required
                            disabled={retrospective.phase > 1 && !isOwner}
                            />
                        <button type="submit" class="hidden" />
                    </form>
                </div>
            </div>
            <div>
                {#each retrospective.workedItems as item(item.id)}
                    <div class="py-1 my-1 item-list-item" data-itemType="worked" data-itemId="{item.id}">
                        <div class="flex" data-dragdisabled={item.items.length > 0}>
                            <div class="flex-shrink">
                                {#if retrospective.phase === 1 || isOwner}
                                    <button on:click={handleItemDelete('worked', item.id)} class="pr-2 pt-1 text-gray-500 hover:text-red-500"><CrossCircle height="18" width="18" /></button>
                                {/if}
                            </div>
                            <div class="flex-grow">
                                <div>{item.content}</div>
                                {#each item.items as child(child.id)}
                                    <div class="flex items-center pl-2 border-l border-gray-300">
                                        <div class="flex-shrink">
                                            {#if retrospective.phase > 1 && isOwner}
                                                <button on:click={unnestItem('worked', child.id)} class="pr-1 text-gray-500 hover:text-green-500"><ArrowLeft /></button>
                                            {/if}
                                        </div>
                                        <div class="flex-grow">{child.content}</div>
                                    </div>
                                {/each}
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        </div>
        <div class="w-1/4 mx-2 p-4 bg-white shadow">
            <div class="flex items-center mb-2">
                <div class="flex-shrink pr-1">
                    <FrownCircle class="text-gray-400" height="24" width="24" />
                </div>
                <div class="flex-grow">
                    <form on:submit={handleNeedsImprovement}>
                        <input
                            bind:value="{needsImprovement}"
                            placeholder="What needs improvement..."
                            class="border-gray-300 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-orange-500"
                            id="needsImprovement"
                            name="needsImprovement"
                            type="text"
                            required
                            disabled={retrospective.phase > 1 && !isOwner} />
                        <button type="submit" class="hidden" />
                    </form>
                </div>
            </div>
            <div>
                {#each retrospective.improveItems as item}
                    <div class="py-1 my-1 item-list-item" data-itemType="improve" data-itemId="{item.id}">
                        <div class="flex" data-dragdisabled={item.items.length > 0}>
                            <div class="flex-shrink">
                                {#if retrospective.phase === 1 || isOwner}
                                    <button on:click={handleItemDelete('improve', item.id)} class="pr-2 pt-1 text-gray-500 hover:text-red-500"><CrossCircle height="18" width="18" /></button>
                                {/if}
                            </div>
                            <div class="flex-grow">
                                <div>{item.content}</div>
                                {#each item.items as child(child.id)}
                                    <div class="flex items-center pl-2 border-l border-gray-300">
                                        <div class="flex-shrink">
                                            {#if retrospective.phase > 1 && isOwner}
                                                <button on:click={unnestItem('improve', child.id)} class="pr-1 text-gray-500 hover:text-green-500"><ArrowLeft /></button>
                                            {/if}
                                        </div>
                                        <div class="flex-grow">{child.content}</div>
                                    </div>
                                {/each}
                            </div>
                        </div>
                    </div>
                {/each}
            </div>
        </div>
        <div class="w-1/4 mx-2 p-4 bg-white shadow">
            <div class="flex items-center mb-2">
                <div class="flex-shrink pr-1">
                    <QuestionCircle class="text-gray-400" height="24" width="24" />
                </div>
                <div class="flex-grow">
                    <form on:submit={handleQuestion}>
                        <input
                            bind:value="{question}"
                            placeholder="I want to ask..."
                            class="border-gray-300 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-orange-500"
                            id="question"
                            name="question"
                            type="text"
                            required
                            disabled={retrospective.phase > 1 && !isOwner}
                            />
                        <button type="submit" class="hidden" />
                    </form>
                </div>
            </div>
            <div>
                {#each retrospective.questionItems as item}
                <div class="py-1 my-1 item-list-item" data-itemType="question" data-itemId="{item.id}">
                    <div class="flex" data-dragdisabled={item.items.length > 0}>
                        <div class="flex-shrink">
                            {#if retrospective.phase === 1 || isOwner}
                                <button on:click={handleItemDelete('question', item.id)} class="pr-2 pt-1 text-gray-500 hover:text-red-500"><CrossCircle height="18" width="18" /></button>
                            {/if}
                        </div>
                        <div class="flex-grow">
                            <div>{item.content}</div>
                            {#each item.items as child(child.id)}
                            <div class="flex items-center pl-2 border-l border-gray-300">
                                <div class="flex-shrink">
                                    {#if retrospective.phase > 1 && isOwner}
                                        <button on:click={unnestItem('question', child.id)} class="pr-1 text-gray-500 hover:text-green-500"><ArrowLeft /></button>
                                    {/if}
                                </div>
                                <div class="flex-grow">{child.content}</div>
                            </div>
                            {/each}
                        </div>
                    </div>
                </div>
                {/each}
            </div>
        </div>
        <div class="w-1/4 mx-2 p-4 bg-white shadow">
            <div class="flex items-center mb-2">
                <div class="flex-shrink pr-1">
                    <CheckCircle class="text-gray-400" height="24" width="24" />
                </div>
                <div class="flex-grow">
                    <form on:submit={handleActionItem}>
                        <input
                            bind:value="{actionItem}"
                            placeholder="Action item..."
                            class="border-gray-300 border-2
                            appearance-none rounded w-full py-2 px-3
                            text-gray-700 leading-tight focus:outline-none
                            focus:bg-white focus:border-orange-500"
                            id="actionItem"
                            name="actionItem"
                            type="text"
                            required
                            disabled={retrospective.phase !== 3}
                            />
                        <button type="submit" class="hidden" />
                    </form>
                </div>
            </div>
            {#each retrospective.actionItems as item, i}
                <div class="py-1 my-1">
                    <div class="flex content-center">
                        <div class="flex-shrink">
                            {#if isOwner}
                                <button on:click={handleActionDelete(item.id)} class="pr-2 pt-1 text-gray-500 hover:text-red-500"><CrossCircle height="18" width="18" /></button>
                            {/if}
                        </div>
                        <div class="flex-grow">
                            {item.content}
                        </div>
                        <div class="flex-shrink">
                            <input type="checkbox" id="{i}Completed" checked="{item.completed}" class="opacity-0 absolute h-6 w-6" on:change={handleActionUpdate(item.id, item.completed)} />
                            <div class="bg-white border-2 rounded-md border-gray-400 w-6 h-6 flex flex-shrink-0 justify-center items-center mr-2 focus-within:border-blue-500">
                                <CheckboxIcon />
                            </div>
                            <label for="{i}Completed" class="select-none"></label>
                        </div>
                    </div>
                </div>
            {/each}
        </div>
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