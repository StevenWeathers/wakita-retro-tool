<script>
    import Sockette from 'sockette'
    import { onMount, onDestroy } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import UserCard from '../components/UserCard.svelte'
    import InviteUser from '../components/InviteUser.svelte'
    import UsersIcon from '../components/icons/UsersIcon.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import DownCarrotIcon from '../components/icons/DownCarrotIcon.svelte'
    import DeleteRetrospective from '../components/DeleteRetrospective.svelte'
    import { appRoutes, PathPrefix } from '../config'
    import { user } from '../stores.js'

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
        users: [],
    }
    let showUsers = false
    let showDeleteRetrospective = false

    const onSocketMessage = function(evt) {
        const parsedEvent = JSON.parse(evt.data)

        switch (parsedEvent.type) {
            case 'init':
                retrospective = JSON.parse(parsedEvent.value)
                eventTag('join', 'retrospective', '')
                break
            case 'user_joined':
                retrospective.users = JSON.parse(parsedEvent.value)
                const joinedUser = retrospective.users.find(
                    w => w.id === parsedEvent.userId,
                )
                notifications.success(`${joinedUser.name} joined.`)
                break
            case 'user_retreated':
                const leftUser = retrospective.users.find(
                    w => w.id === parsedEvent.userId,
                )
                retrospective.users = JSON.parse(parsedEvent.value)

                notifications.danger(`${leftUser.name} retreated.`)
                break
            case 'retrospective_updated':
                retrospective = JSON.parse(parsedEvent.value)
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

    onMount(() => {
        if (!$user.id) {
            router.route(`${appRoutes.login}/${retrospectiveId}`)
        }
    })
</script>

<svelte:head>
    <title>Retrospective {retrospective.name} | Wakita</title>
</svelte:head>

{#if retrospective.name && !socketReconnecting && !socketError}
    <div class="px-6 py-2 bg-gray-200 flex flex-wrap">
        <div class="w-1/3">
            <h1 class="text-3xl font-bold leading-tight">
                {retrospective.name}
            </h1>
        </div>
        <div class="w-2/3 text-right">
            <div>
                {#if retrospective.ownerId === $user.id}
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
                                            {sendSocketEvent}
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
    <div>Retrospective board here...</div>
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