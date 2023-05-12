import type { App } from 'vue';
import type { Directive } from 'vue';

interface ExtendedHTMLElement extends HTMLElement {
  clickOutsideEvent(target: Event): void | null
}

export default {
  install: (app: App) => {
    app.directive('clickOutside', <Directive<HTMLElement, string>>{
      beforeMount(el: ExtendedHTMLElement, binding, vnode, prevVnode) {
        if (typeof binding.value !== 'function') {
          console.warn(`[v-click-outside] binding value '${binding.value}' is not a function`);
          return;
        }
        el.clickOutsideEvent = (event: MouseEvent) => {
          if (!el.contains(<Node>event.target) && el !== event.target) {
            (<any>binding.value)(event);
          }
        };
        document.addEventListener('click', el.clickOutsideEvent);
      },
      beforeUnmount(el: ExtendedHTMLElement, binding, vnode, prevVnode) {
        document.removeEventListener('click', el.clickOutsideEvent)
      },
    })
  }
};