import { defineStore } from 'pinia'
import axios from 'axios'

export const useTimeStore = defineStore('time', {
    state: () => ({
        now: null,
        range: 1800,
        offset: new Date().getTimezoneOffset(),
        status: { ready: false, loading: false, retries: 0 },
        timerIDs: [],
        timerManager: '',
    }),
    computed: {
        nowUnix() {
            return Date.parse(this.nowDate)
        }
    },
    actions: {
        async toTimeRangeForQuery(r, updateNow = true) {
            const now = await this.getNow(updateNow)
            return [await this.toAbsoluteTime(r[0], now), await this.toAbsoluteTime(r[1], now)]
        },
        async toAbsoluteTime(t, now) {
            if (!t.startsWith('now')) {
                return Date.parse(t) / 1000
            }
            if (t == 'now') return now
            const offset = t.split('-').pop()
            let num = offset.slice(0, -1)
            switch (offset.slice(-1)) {
                case 'm': num *= 60; break;
                case 'h': num *= 3600; break;
                case 'd': num *= 86400; break;
            }
            return now - num
        },
        async getNow(updateNow = true) {
            if (updateNow || !this.now) await this.updateNow()
            return this.now
        },
        async updateNow() {
            try {
                const response = await axios.get('/api/prometheus/time')
                const now = response.data.data.result[0]
                this.now = now
            } catch (error) {
                console.error(error)
            }
        },
        getOffset() {
            return this.offset
        },
        getOffsetString() {
            const o = this.offset
            const sign = (o > 0) ? '-' : '+'
            const hours = ('0' + Math.abs(o) / 60).slice(-3)
            const minutes = ('0' + Math.abs(o) % 60).slice(-3)
            return `UTC${sign}${hours}:${minutes}`
        },
        timestamp2ymdhis(t) {
            return new Date(t * 1000 - this.offset * 60000).toISOString().replace(/\..*/, '').replace(/T/, ' ')
        },
    }
})
