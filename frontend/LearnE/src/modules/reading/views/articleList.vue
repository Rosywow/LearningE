<template>
    <h1>文章列表</h1>
    <ul>
        <li v-for="(title, index) in titles " :key="index">
            <router-link :to="{ name: 'article', params: { id: title } }">{{ title }}</router-link>
        </li>
    </ul>
</template>

<script setup>
import { onMounted, ref } from "vue";
onMounted(() => {
    let url = `/api/reading/getArticleList`
    fetch(url, {
        method: "get"
    })
        .then(res => res.json())
        .then(data => {
            titles.value = data.data
            console.log(titles.value)
        })
        .catch(err => console.log("err:", err))
})
const titles = ref([])
</script>
