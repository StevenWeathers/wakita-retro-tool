<script>
    import SmileCircle from './icons/SmileCircle.svelte'
    import CrossCircle from './icons/CrossCircle.svelte'
    import FrownCircle from './icons/FrownCircle.svelte'
    import QuestionCircle from './icons/QuestionCircle.svelte'
    import ThumbsUp from './icons/ThumbsUp.svelte'
    import ArrowLeft from './icons/ArrowLeft.svelte'

    export let handleSubmit = () => {}
    export let handleDelete = () => {}
    export let handleVote = () => {}
    export let handleUnnest = () => {}
    export let itemType = 'worked'
    export let content = ''
    export let newItemPlaceholder = 'What worked well...'
    export let phase = 1
    export let isOwner = false
    export let items = []

    const handleFormSubmit = (evt) => {
        evt.preventDefault()

        handleSubmit(itemType, content)
        content = ''
    }
</script>

<div class="w-1/4 mx-2 p-4 bg-white shadow">
    <div class="flex items-center mb-2">
        <div class="flex-shrink pr-1">
            {#if itemType === 'worked'}
                <SmileCircle class="text-gray-400" height="24" width="24" />
            {:else if itemType === 'improve'}
                <FrownCircle class="text-gray-400" height="24" width="24" />
            {:else if itemType === 'question'}
                <QuestionCircle class="text-gray-400" height="24" width="24" />
            {/if}
        </div>
        <div class="flex-grow">
            <form on:submit={handleFormSubmit} class="flex">
                <input
                    bind:value="{content}"
                    placeholder="{newItemPlaceholder}"
                    class="border-gray-300 border-2
                    appearance-none rounded py-2 px-3
                    text-gray-700 leading-tight focus:outline-none
                    focus:bg-white focus:border-orange-500 w-full"
                    id="new{itemType}"
                    name="new{itemType}"
                    type="text"
                    required
                    disabled={phase > 1 && !isOwner}
                    />
                <button type="submit" class="hidden" />
            </form>
        </div>
    </div>
    <div>
        {#each items as item(item.id)}
            <div class="py-1 my-1 item-list-item" data-itemType="{itemType}" data-itemId="{item.id}">
                <div class="flex" data-dragdisabled={item.items.length > 0}>
                    <div class="flex-shrink">
                        {#if phase === 1 || isOwner}
                            <button on:click={handleDelete(itemType, item.id)} class="pr-2 pt-1 text-gray-500 hover:text-red-500"><CrossCircle height="18" width="18" /></button>
                        {/if}
                    </div>
                    <div class="flex-grow">
                        <div class="flex items-center">
                            <div class="flex-grow">{item.content}</div>
                            <div class="flex-shrink">
                                {#if phase > 1}
                                    <button on:click={handleVote(itemType, item.id)} class="pr-1 text-gray-500 hover:text-green-500" disabled={phase !== 2}><ThumbsUp /></button>
                                    <span class="text-gray-600">&nbsp;{item.votes.length}</span>
                                {/if}
                            </div>
                        </div>
                        {#each item.items as child(child.id)}
                            <div class="flex items-center pl-2 pt-1 border-l border-gray-300">
                                <div class="flex-shrink">
                                    {#if phase > 1 && isOwner}
                                        <button on:click={handleUnnest(itemType, child.id)} class="pr-1 text-gray-500 hover:text-green-500"><ArrowLeft /></button>
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