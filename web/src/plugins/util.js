import ColorHash from 'color-hash'
const colorHash = new ColorHash({ lightness: [0.4, 0.55, 0.7], saturation: [0.4, 0.55, 0.7] })
export default {
    install: ({ config }) => {
        config.globalProperties.$util = {
            dateTimeAsLocal(d) {
                d = new Date(parseFloat(d * 1000))
                d = new Date(d.getTime() - (d.getTimezoneOffset() * 60 * 1000))
                return d.toISOString().substr(0, 19).replace('T', ' ')
            },
            utc2age(u) {
                return this.seconds2age((new Date() - new Date(u)) / 1000)
            },
            age2seconds(a) {
                if (!a && a.length < 2) return 0
                let num = a.slice(0, -1)
                switch (a.slice(-1)) {
                    case 's': return num * 1
                    case 'm': return num * 60
                    case 'h': return num * 3600
                    case 'd': return num * 86400
                }
                return 0
            },
            seconds2age(s) {
                if (s == 0) return 0
                if (s > 86400) return Math.floor(s / 86400) + "d"
                if (s > 3600) return Math.floor(s / 3600) + "h"
                if (s > 60) return Math.floor(s / 60) + "m"
                return Math.floor(s) + "s"
            },
            nanoseconds2human(ns) {
                if(!ns)ns=0
                return this.seconds2age(ns/1000000000)
            },
            string2color(str) {
                return colorHash.hex(str)
            },
            string2letters(str) {
                return str.split(/[.-]+/).map(x => x[0]).join('')
            },
            copyToClipboard(val) {
                const t = document.createElement("textarea")
                document.body.appendChild(t)
                t.value = val
                t.select()
                document.execCommand('copy') // deprecated
                document.body.removeChild(t)
            },
        }
    }
}