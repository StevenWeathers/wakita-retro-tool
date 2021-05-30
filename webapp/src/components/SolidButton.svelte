<script>
    export let disabled = false
    export let color = 'green'
    export let additionalClasses = ''
    export let type = 'button'
    export let onClick = () => {}
</script>

<style>
    .btn {
        @apply leading-tight;
        @apply py-2;
        @apply px-3;
        @apply rounded;
        @apply font-bold;
        @apply text-white;
        @apply border;
    }
    .btn.disabled {
        @apply opacity-50;
        @apply cursor-not-allowed;
    }
    .btn-green {
        @apply bg-green-500;
        @apply border-green-500;
    }
    .btn-green:hover {
        @apply bg-green-600;
        @apply border-green-600;
    }
    .btn-blue {
        @apply bg-blue-500;
        @apply border-blue-500;
    }
    .btn-blue:hover {
        @apply bg-blue-600;
        @apply border-blue-600;
    }
    .btn-red {
        @apply bg-red-500;
        @apply border-red-500;
    }
    .btn-red:hover {
        @apply bg-red-600;
        @apply border-red-600;
    }
    .btn-purple {
        @apply bg-purple-500;
        @apply border-purple-500;
    }
    .btn-purple:hover {
        @apply bg-purple-600;
        @apply border-purple-600;
    }
    .btn-teal {
        @apply bg-teal-500;
        @apply border-teal-500;
    }
    .btn-teal:hover {
        @apply bg-teal-600;
        @apply border-teal-600;
    }
    .btn-blue-copy {
        @apply bg-blue-500;
        @apply rounded-l-none;
        @apply border-blue-800;
    }
    .btn-blue-copy:hover {
        @apply bg-blue-600;
        @apply border-blue-600;
    }
    .btn-yellow {
        @apply bg-yellow-500;
        @apply border-yellow-500;
    }
    .btn-yellow:hover {
        @apply bg-yellow-600;
        @apply border-yellow-600;
    }
    .btn-orange {
        @apply bg-orange-500;
        @apply border-orange-500;
    }
    .btn-orange:hover {
        @apply bg-orange-600;
        @apply border-orange-600;
    }
</style>

<button
    class="btn btn-{color}
    {disabled ? 'disabled' : ''}
    {additionalClasses}"
    on:click="{onClick}"
    {type}
    {disabled}>
    <slot />
</button>
