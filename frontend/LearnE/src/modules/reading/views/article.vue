<template>
  <div>
    <h2>{{ title }}</h2>
  </div>
  <p>
    <span v-for="(word, index) in words" :key="index">
      <a class="word" :class="{ selectedWord: word === clickedWords }" @click="storeWord(word)">{{ word + " "
      }}</a></span>
  </p>
</template>

<script setup>
import axios from "axios"
import md5 from "md5"
import { useRoute } from "vue-router";
import { ref, onMounted } from "vue";
const text = ref("");
const title = ref("");
const article = ref("");
const words = ref("");
const clickedWords = ref("");
const route = useRoute();
onMounted(() => {
  title.value = route.params.id;
  let url = `/api/reading/getArticle?title=${title.value}`;
  fetch(url, {
    method: "post",
  })
    .then((res) => res.json())
    .then((data) => {
      article.value = data;
      words.value = article.value.split(/(\s+|\b)/);;
    })
    .catch((err) => console.log("err", err));
});
const storeWord = async (word) => {
  let url = `/api/reading/storeWord`;
  var utterance = new SpeechSynthesisUtterance();
  utterance.text = word;
  window.speechSynthesis.speak(utterance);
  if (clickedWords.value != word) {
    await fetch(url, {
      method: "POST",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({
        word: word,
      }),
    })
      .then(() => {
        clickedWords.value = word
      })
      .catch((err) => console.log(err));
  }
};
</script>

<style>
.word {
  /* 保留空格，换行符，制表符，可以自动换行，将多个连续普通空格合并为一个空格 */
  white-space: pre-line;
}

.selectedWord {
  background-color: aqua;
}
</style>
