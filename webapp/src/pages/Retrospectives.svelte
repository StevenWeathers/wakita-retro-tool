<script>
    import { onMount } from 'svelte'

    import PageLayout from '../components/PageLayout.svelte'
    import CreateRetrospective from '../components/CreateRetrospective.svelte'
    import HollowButton from '../components/HollowButton.svelte'
    import { user } from '../stores.js'
    import { appRoutes } from '../config'

    export let xfetch
    export let notifications
    export let router
    export let eventTag

    let retrospectives = []

    xfetch('/api/retrospectives')
        .then(res => res.json())
        .then(function(bs) {
            retrospectives = bs
        })
        .catch(function(error) {
            notifications.danger('Error finding your retrospectives')
            eventTag('fetch_retrospectives', 'engagement', 'failure')
        })

    onMount(() => {
        if (!$user.id) {
            router.route(appRoutes.login)
        }
    })
</script>

<PageLayout>
    <h1 class="mb-4 text-3xl font-bold">My Retrospectives</h1>

    <div class="flex flex-wrap">
        <div class="mb-4 md:mb-6 w-full md:w-1/2 lg:w-3/5 md:pr-4">
            {#each retrospectives as retrospective}
                <div class="bg-white shadow-lg rounded mb-2">
                    <div
                        class="flex flex-wrap items-center p-4 border-gray-400
                        border-b">
                        <div
                            class="w-full md:w-1/2 mb-4 md:mb-0 font-semibold
                            md:text-xl leading-tight">
                            {retrospective.name}
                            <div class="font-semibold md:text-sm text-gray-600">
                                {#if $user.id === retrospective.owner_id}Owner{/if}
                            </div>
                        </div>
                        <div class="w-full md:w-1/2 md:mb-0 md:text-right">
                            <HollowButton
                                href="{appRoutes.retrospective}/{retrospective.id}">
                                Join Retro
                            </HollowButton>
                        </div>
                    </div>
                </div>
            {/each}
        </div>

        <div class="w-full md:w-1/2 lg:w-2/5 pl-4">
            <div class="p-6 bg-white shadow-lg rounded">
                <h2 class="mb-4 text-2xl font-bold leading-tight">
                    Create a Retro
                </h2>
                <CreateRetrospective
                    {notifications}
                    {router}
                    {eventTag}
                    {xfetch} />
            </div>
        </div>
    </div>
</PageLayout>
