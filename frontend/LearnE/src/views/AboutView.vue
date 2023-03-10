<template>
  <h1>space repetition - english learning</h1>
  <ul>
    <li><router-link :to="{ name: 'addsource' }">addsource</router-link></li>
    <li><router-link :to="{ name: 'articleList' }">articleList</router-link></li>
    <li><router-link :to="{ name: 'video' }">video</router-link></li>
  </ul>
</template>

<script setup lang="ts">
import { onMounted } from "vue"
onMounted(() => {

})
</script>

<style>
ul {
  list-style: decimal;
}
</style>
