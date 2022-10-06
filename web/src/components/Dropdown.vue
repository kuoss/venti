<template>
    <div class="inline-block relative z-30">
        <button class="border text-gray-700 py-2 px-4 rounded inline-flex items-center" @click.stop="toggleIsOpen">
            <span v-if="selected">{{ selected }}</span>
            <span v-else>{{ options[0] }}</span>
            <i class="mdi mdi-chevron-down"></i>
        </button>
        <ul v-if="isOpen" class="absolute text-gray-700 w-max border border-gray-300 bg-white">
            <li
                class="hover:bg-gray-400 px-2 cursor-pointer py-[1px]"
                @click.stop="select(option)"
                v-for="option in options"
            >{{ option }}</li>
        </ul>
    </div>
</template>

<script>
let uid = 0
export default {
    props: {
        options: Array,
        currentDropdown: Number,
    },
    watch: {
        currentDropdown(newCurrentDropdown) {
            if(this.uid != newCurrentDropdown) this.close()
        },
    },
    data() {
        return {
            isOpen: false,
            selected: null,
        }
    },
    methods: {
        toggleIsOpen() {
            this.isOpen = !this.isOpen
            if (this.isOpen) {
                this.$emit('open', this.uid)
            }
        },
        select(option) {
            this.isOpen = false
            this.selected = option
            this.$emit('select', option)
        },
        close() {
            this.isOpen = false
        },
        clickOutside(event, object) {
            console.log(event, object)
            this.close()
        },
    },
    beforeCreate() {
        this.uid = ++uid
    },
    mounted() {
        document.addEventListener("click", this.close)
    },
    unmounted() {
        document.removeEventListener("click", this.close)
    },
}
</script>