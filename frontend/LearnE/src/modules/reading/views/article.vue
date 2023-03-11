<template>
  <div>
    <h2>{{ title }}</h2>
  </div>
  <a-button @click="countDown">Open modal to close in 5s</a-button>
  <p>
    <span v-for="(word, index) in words" :key="index">
      <a class="word" :class="{ selectedWord: word === clickedWords }" @click="storeWord(word)">{{ word + " "
      }}</a></span>
  </p>
</template>

<script setup>
import { Modal } from 'ant-design-vue';
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
const showWord = (word, meaning) => {
  // 2秒后消失
  let secondsToGo = 2;
  const modal = Modal.success({
    title: word,
    content: meaning,
  });
  const interval = setInterval(() => {
    secondsToGo -= 1;
  }, 1000);
  setTimeout(() => {
    clearInterval(interval);
    modal.destroy();
  }, secondsToGo * 1000);
};
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
  // let url = `/api/reading/storeWord`;
  // var utterance = new SpeechSynthesisUtterance();
  // utterance.text = word;
  // window.speechSynthesis.speak(utterance);
  // if (clickedWords.value != word) {
  //   await fetch(url, {
  //     method: "POST",
  //     headers: {
  //       "Content-Type": "application/json",
  //     },
  //     body: JSON.stringify({
  //       word: word,
  //     }),
  //   })
  //     .then(() => {
  //       clickedWords.value = word
  //     })
  //     .catch((err) => console.log(err));
  // }
  let salt = Math.random()
  axios.get(`/trans`, {
    params: {
      q: word,
      from: 'en',
      to: "zh",
      appid: "20230304001585991",
      salt: salt,
      sign: md5("20230304001585991" + word + salt + "UugK4LkE7YTFbhEZoEnO")
    }
  })
    .then(response => {
      console.log(response.data.trans_result[0].dst)
      showWord(word, response.data.trans_result[0].dst)
    })
    .catch(err => {
      console.log(err)
    })
  console.log(word)
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
