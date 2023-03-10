<template>
  <form @submit.prevent="upLoad()">
    <div class="container">
      <input placeholder="title" v-model="title" required :style="{ width: '250px', textalign: 'right' }">
    </div>
    <div class="container">
      <textarea v-model="article" required :style="{ width: '600px', height: '450px', textalign: 'right' }"></textarea>
    </div>
    <div class="container">
      <input placeholder="Link" v-model="link"  :style="{ width: '300px', textalign: 'right' }">
    </div>
    <div class="container">
      <button type="submit">上传</button>
    </div>
  </form>
</template>

<script setup>
import { ref } from 'vue'
const title = ref("")
const article = ref("")
const link = ref("")
const upLoad = () => {
  let url = `/api/reading/addsource`;
  fetch(url, {
    method: "POST",
    headers: {
      'Content-Type': 'application/json'
    },
    body: JSON.stringify({
      title: title.value,
      article: article.value,
      link: link.value,
    })
  })
    .then(res => res.json())
    .then(data => {
      title.value = "";
      article.value = ""
      link.value = ""
    })
    .catch(err => console.log(err))
}
</script>

<style scoped>
.container {
  display: flex;
  align-items: center;
  justify-content: center;
  margin-bottom: 20px;
}
</style>