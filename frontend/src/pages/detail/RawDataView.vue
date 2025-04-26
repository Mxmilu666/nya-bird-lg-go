<template>
    <div class="detail-section">
        <h3>Raw Output</h3>
        <a-textarea :rows="20" :value="rawOutput" readonly />
    </div>
</template>

<script setup lang="ts">
defineProps<{
    rawOutput: string
}>()
</script>

<style scoped/>
