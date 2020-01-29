<script>
    import {createEventDispatcher} from 'svelte';
    const dispatch = createEventDispatcher();

    let text = ''
    let sendMessage = () => {
        dispatch('message', {
            text: 'hello',
        })
    }
</script>

text: <input bind:value={text} />
<button on:cilck={sendMessage}>send</button>
