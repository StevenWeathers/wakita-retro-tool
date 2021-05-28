<script>
    import { onMount } from 'svelte'

    import AdminPageLayout from '../../components/AdminPageLayout.svelte'
    import HollowButton from '../../components/HollowButton.svelte'
    import { user } from '../../stores.js'
    import { _ } from '../../i18n'
    import { appRoutes } from '../../config'

    export let xfetch
    export let router
    export let notifications
    export let eventTag

    const {
        CleanupGuestsDaysOld,
        CleanupRetrospectivesDaysOld,
        APIEnabled,
    } = appConfig

    let appStats = {
        unregisteredUserCount: 0,
        registeredUserCount: 0,
        retrospectiveCount: 0,
        organizationCount: 0,
        departmentCount: 0,
        teamCount: 0,
        apikeyCount: 0,
        activeRetroCount: 0,
        activeRetroUserCount: 0,
    }

    function getAppStats() {
        xfetch('/api/admin/stats')
            .then(res => res.json())
            .then(function(result) {
                appStats = result
            })
            .catch(function(error) {
                notifications.danger('Error getting application stats')
            })
    }

    function cleanRetrospectives() {
        xfetch('/api/admin/clean-retrospectives', { method: 'DELETE' })
            .then(function() {
                eventTag('admin_clean_retrospectives', 'engagement', 'success')

                getAppStats()
            })
            .catch(function(error) {
                notifications.danger(
                    'Error encountered cleaning retrospectives',
                )
                eventTag('admin_clean_retrospectives', 'engagement', 'failure')
            })
    }

    function cleanGuests() {
        xfetch('/api/admin/clean-guests', { method: 'DELETE' })
            .then(function() {
                eventTag('admin_clean_guests', 'engagement', 'success')

                getAppStats()
            })
            .catch(function(error) {
                notifications.danger('Error encountered cleaning guests')
                eventTag('admin_clean_guests', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$user.id) {
            router.route(appRoutes.login)
        }
        if ($user.type !== 'ADMIN') {
            router.route(appRoutes.landing)
        }

        getAppStats()
    })
</script>

<AdminPageLayout activePage="admin">
    <div class="text-center px-2 mb-4">
        <h1 class="text-3xl md:text-4xl font-bold">Admin</h1>
    </div>
    <div class="flex justify-center mb-4">
        <div class="w-full">
            <div
                class="flex flex-wrap items-center text-center pt-2 pb-2 md:pt-4
                md:pb-4 bg-white shadow-lg rounded text-xl">
                <div class="w-1/2">
                    <div class="mb-2 font-bold">Active Retrospectives</div>
                    {appStats.activeRetroCount}
                </div>
                <div class="w-1/2">
                    <div class="mb-2 font-bold">Active Retrospective Users</div>
                    {appStats.activeRetroUserCount}
                </div>
            </div>
            <div
                class="flex flex-wrap items-center text-center pt-2 pb-2 md:pt-4
                md:pb-4 bg-white shadow-lg rounded text-xl">
                <div class="w-1/3">
                    <div class="mb-2 font-bold">Unregistered Users</div>
                    {appStats.unregisteredUserCount}
                </div>
                <div class="w-1/3">
                    <div class="mb-2 font-bold">Registered Users</div>
                    {appStats.registeredUserCount}
                </div>
                <div class="w-1/3">
                    <div class="mb-2 font-bold">Retrospectives</div>
                    {appStats.retrospectiveCount}
                </div>
            </div>
            <div
                class="flex flex-wrap items-center text-center pt-2 pb-2 md:pt-4
                md:pb-4 bg-white shadow-lg rounded text-xl">
                <div class="{APIEnabled ? 'w-1/4' : 'w-1/3'}">
                    <div class="mb-2 font-bold">Organizations</div>
                    {appStats.organizationCount}
                </div>
                <div class="{APIEnabled ? 'w-1/4' : 'w-1/3'}">
                    <div class="mb-2 font-bold">Departments</div>
                    {appStats.departmentCount}
                </div>
                <div class="{APIEnabled ? 'w-1/4' : 'w-1/3'}">
                    <div class="mb-2 font-bold">Teams</div>
                    {appStats.teamCount}
                </div>
                {#if APIEnabled}
                    <div class="w-1/4">
                        <div class="mb-2 font-bold">API Keys</div>
                        {appStats.apikeyCount}
                    </div>
                {/if}
            </div>
        </div>
    </div>

    <div class="flex justify-center mb-4">
        <div class="w-full">
            <div
                class="text-center p-2 md:p-4 bg-white shadow-lg rounded text-xl">
                <div class="text-2xl md:text-3xl font-bold text-center mb-4">
                    Maintenance
                </div>
                <HollowButton onClick="{cleanGuests}" color="red">
                    Clean Guests older than {CleanupGuestsDaysOld} days
                </HollowButton>

                <HollowButton onClick="{cleanRetrospectives}" color="red">
                    Clean Retrospectives older than {CleanupRetrospectivesDaysOld}
                    days
                </HollowButton>
            </div>
        </div>
    </div>
</AdminPageLayout>
