<script>
    import { onMount } from 'svelte'

    import SolidButton from './SolidButton.svelte'
    import { user } from '../stores.js'
    import { appRoutes } from '../config'

    export let xfetch
    export let notifications
    export let eventTag
    export let router
    export let apiPrefix = '/api'

    let retrospectiveName = ''

    function createRetrospective(e) {
        e.preventDefault()
        const body = {
            retrospectiveName,
        }

        xfetch(`${apiPrefix}/retrospective`, { body })
            .then(res => res.json())
            .then(function(retrospective) {
                eventTag('create_retrospective', 'engagement', 'success', () => {
                    router.route(`${appRoutes.retrospective}/${retrospective.id}`)
                })
            })
            .catch(function(error) {
                notifications.danger('Error encountered creating retrospective')
                eventTag('create_retrospective', 'engagement', 'failure')
            })
    }

    onMount(() => {
        if (!$user.id) {
            router.route(appRoutes.register)
        }
    })
</script>

<form on:submit="{createRetrospective}" name="createRetrospective">
    <div class="mb-4">
        <label
            class="block text-gray-700 text-sm font-bold mb-2"
            for="retrospectiveName">
            Retrospective Name
        </label>
        <div class="control">
            <input
                name="retrospectiveName"
                bind:value="{retrospectiveName}"
                placeholder="Enter a retrospective name"
                class="bg-gray-200 border-gray-200 border-2 appearance-none
                rounded w-full py-2 px-3 text-gray-700 leading-tight
                focus:outline-none focus:bg-white focus:border-orange-500"
                id="retrospectiveName"
                required />
        </div>
    </div>

    <div class="text-right">
        <SolidButton type="submit">Create Retrospective</SolidButton>
    </div>
</form>
