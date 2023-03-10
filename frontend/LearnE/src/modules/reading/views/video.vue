<template>
    <h1>video</h1>
    <video class="video" ref="videogame" :style="{ width: '600px' }">
        <source src="../video/videoplayback.mp4" type="video/mp4">
    </video>
    <div class="progress-bar">
        <div class="progress" ref="progress"></div>
    </div>
    <div>{{ lyric }}</div>
    <div>{{ lyric1 }}</div>
    <button @click="playAudio" :style="{ width: '60px', height: '40px' }">{{ state }}</button>
    <div v-for="(sentence, index) in lyrics" :key="index" :class="{ selectedSentence: index === currentIndex }">
        <span v-for="(word, i) in sentence.text.split(/(\s+|\b)/)" :key="i" @click="StoreWord(word)">
            {{ word }}
        </span>
    </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, reactive, ref, watch } from 'vue'
import { parse } from 'lrc-parser'
const state = ref("▶")
const videogame = ref(null)
const progress = ref(null)
const text = ref("")
type Lyric = { time: number; text: string }
const lyrics = reactive<Lyric[]>([])
const lyric = ref("")
const lyric1 = ref("")
const currentIndex = ref(0)
const clickedWord = ref("")
let timer;
const parseLRC = (context: string) => {
    const lines = context.split('\n')
    var timeReg = /\[(\d{1,2}):(\d{1,2})(?::(\d{1,2}))?]/;

    for (var i = 0; i < lines.length; i++) {
        var matches = lines[i].match(timeReg);
        let time;
        if (matches != null && matches[3]) {
            // 说明有小时 时间戳格式为xx:xx:xx
            let hour = parseInt(matches[1])
            let minutes = parseInt(matches[2])
            let seconds = parseInt(matches[3])
            console.log(hour, minutes, seconds)
            time = hour * 3600 + minutes * 60 + seconds
        } else if (matches != null && !matches[3]) {
            // xx:xx
            let minutes = parseInt(matches[1]);
            let seconds = parseInt(matches[2]);
            time = minutes * 60 + seconds
        }
        // 识别文字
        var text = lines[i].replace(timeReg, '').trim();
        lyrics.push({
            time: time,
            text: text,
        })
    }
}
onMounted(() => {
    window.addEventListener('keydown', e => {
        e.preventDefault(); // 阻止默认行为
        if (e.key === " ") {
            playAudio()
        }
    })
    // 定时同步字幕
    timer = setInterval(() => {
        // 获取当前时间
        if (state.value !== "▶") {
            //@ts-ignore
            let currentTime = videogame.value.currentTime
            for (var i = 1; i < lyrics.length - 1; i++) {
                if (lyrics[i - 1].time < currentTime && lyrics[i].time > currentTime) {
                    currentIndex.value = i - 1
                    lyric.value = lyrics[i - 1].text
                    lyric1.value = lyrics[i].text
                }
            }
        }
    }, 2000)
    //@ts-ignore
    videogame.value.addEventListener('timeupdate', updateProgressBar);
    // 获取字幕
    fetch(`/api/reading/geiLyric`, {
        method: "GET"
    })
        .then(res => res.json())
        .then(data => {
            text.value = data
            parseLRC(text.value)
        })
        .catch(err => console.log(err))
})
onUnmounted(() => {
    clearInterval(timer)
    window.removeEventListener('keydown', e => {
        e.preventDefault(); // 阻止默认行为
        if (e.key === " ") {
            playAudio()
        }
    })
})
const playAudio = () => {
    console.log(videogame)
    //@ts-ignore
    if (videogame.value.paused) {
        //@ts-ignore
        console.log("播放", videogame.value.play())
        state.value = "⏸"
    } else {
        //@ts-ignore
        console.log("暂停", videogame.value?.pause())
        state.value = "▶"
    }
}
const updateProgressBar = () => {
    //@ts-ignore
    const percent = (videogame.value.currentTime / videogame.value.duration) * 100;
    //@ts-ignore
    progress.value.style.width = percent + '%';
}
const storeWord = async (word) => {
    let url = `/api/reading/storeWord`;
    var utterance = new SpeechSynthesisUtterance();
    utterance.text = word;
    window.speechSynthesis.speak(utterance);
    if (clickedWord.value != word) {
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
                clickedWord.value = word
            })
            .catch((err) => console.log(err));
    }
};
</script>

<style scoped>
.progress-bar {
    width: 600px;
    height: 5px;
    background-color: #ccc;
}

.progress {
    width: 0%;
    height: 100%;
    background-color: #00bcd4;
    transition: width 0.1s linear;
}

.selectedSentence {
    background-color: #00bcd4;
}
</style>