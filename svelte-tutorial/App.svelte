<script>
    import Nested from './Nestes.svelte';
    let getMessage = (m) => {
        alert(m)
        console.log(m)
    };
</script>

<Nested on:message={getMessage}/>

