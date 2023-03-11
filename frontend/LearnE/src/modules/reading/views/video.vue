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
    <div class="sentence" v-for="(sentence, index) in lyrics" :key="index"
        :class="{ selectedSentence: index === currentIndex }" @dblclick="Jump(sentence.time)">{{ sentence.time + ": " }}
        <span v-for="(word, i) in sentence.text.split(/(\s+|\b)/)" :key="i" @click="storeWord(word)">
            {{ word }}
        </span>
    </div>
</template>

<script setup lang="ts">
import { Modal } from 'ant-design-vue';
import { onMounted, onUnmounted, reactive, ref } from 'vue'
import axios from "axios"
import md5 from "md5"
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
// 双击句子实现视频跳转
const Jump = (jumpTo) => {
    // @ts-ignore
    videogame.value.currentTime = jumpTo
    console.log("dbclick")
}
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
        if (e.key === " ") {
            playAudio()
            e.preventDefault(); // 阻止默认行为
        } else if (e.key === "ArrowLeft" && currentIndex.value > 0) {
            //@ts-ignore
            videogame.value.currentTime = lyrics[currentIndex.value - 1].time
            currentIndex.value--
        } else if (e.key === "ArrowRight" && currentIndex.value < lyrics.length - 1) {
            //@ts-ignore
            videogame.value.currentTime = lyrics[currentIndex.value + 1].time
            currentIndex.value++
        } else if (e.key === "Enter") {
            console.log("enter")
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
    console.log(state.value)
    //@ts-ignore
    if (videogame.value.paused) {
        //@ts-ignore
        videogame.value.play()
        state.value = "⏸"
    } else {
        //@ts-ignore
        videogame.value?.pause()
        state.value = "▶"
    }
}
const updateProgressBar = () => {
    //@ts-ignore
    const percent = (videogame.value.currentTime / videogame.value.duration) * 100;
    //@ts-ignore
    progress.value.style.width = percent + '%';
}
const showWord = (word) => {
    let salt = Math.random()
    // 百度翻译api
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
            // 2秒后消失
            let secondsToGo = 2;
            const modal = Modal.success({
                title: word,
                content: response.data.trans_result[0].dst,
            });
            const interval = setInterval(() => {
                secondsToGo -= 1;
            }, 1000);
            setTimeout(() => {
                clearInterval(interval);
                modal.destroy();
            }, secondsToGo * 1000);
        })
        .catch(err => {
            console.log(err)
        })

};
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
                showWord(word)
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

.sentence {
    margin-bottom: 5px;
}
</style>