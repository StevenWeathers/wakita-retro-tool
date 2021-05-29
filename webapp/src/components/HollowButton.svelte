<script>
    export let disabled = false
    export let color = 'green'
    export let additionalClasses = ''
    export let type = 'button'
    export let onClick = () => {}
    export let href = ''
    export let btnSize = ''

    $: btnSizeClass = btnSize === 'small' ? 'btn-hollow-small' : ''
</script>

<style>
    .btn-hollow {
        @apply leading-tight;
        @apply font-semibold;
        @apply bg-transparent;
        @apply py-2;
        @apply px-3;
        @apply border;
        @apply rounded;
    }
    .btn-hollow-small {
        @apply py-1;
        @apply px-2;
    }
    .btn-hollow.disabled {
        @apply opacity-50;
        @apply cursor-not-allowed;
    }
    .btn-hollow:hover {
        @apply border-transparent;
        @apply text-white;
    }
    .btn-hollow-green {
        @apply border-green-500;
        @apply text-green-500;
    }
    .btn-hollow-green:hover {
        @apply bg-green-600;
    }
    .btn-hollow-blue {
        @apply border-blue-500;
        @apply text-blue-500;
    }
    .btn-hollow-blue:hover {
        @apply bg-blue-600;
    }
    .btn-hollow-red {
        @apply border-red-500;
        @apply text-red-500;
    }
    .btn-hollow-red:hover {
        @apply bg-red-600;
    }
    .btn-hollow-purple {
        @apply border-purple-500;
        @apply text-purple-500;
    }
    .btn-hollow-purple:hover {
        @apply bg-purple-600;
    }
    .btn-hollow-teal {
        @apply border-teal-500;
        @apply text-teal-500;
    }
    .btn-hollow-teal:hover {
        @apply bg-teal-600;
    }
    .btn-hollow-yellow {
        @apply border-yellow-500;
        @apply text-yellow-500;
    }
    .btn-hollow-yellow:hover {
        @apply bg-yellow-600;
    }
    .btn-hollow-orange {
        @apply border-orange-500;
        @apply text-orange-500;
    }
    .btn-hollow-orange:hover {
        @apply bg-orange-600;
    }
    .btn-hollow-gray {
        @apply border-gray-500;
        @apply text-gray-500;
    }
    .btn-hollow-gray:hover {
        @apply bg-gray-600;
    }
    .btn-hollow-indigo {
        @apply border-indigo-500;
        @apply text-indigo-500;
    }
    .btn-hollow-indigo:hover {
        @apply bg-indigo-600;
    }
</style>

{#if href === ''}
    <button
        class="btn-hollow {btnSizeClass} btn-hollow-{color}
        {disabled ? 'disabled' : ''}
        {additionalClasses}"
        on:click="{onClick}"
        {type}
        {disabled}>
        <slot />
    </button>
{:else}
    <a
        {href}
        class="btn-hollow btn-hollow-{color} inline-block no-underline {additionalClasses}">
        <slot />
    </a>
{/if}
