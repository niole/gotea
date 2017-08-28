import Vue from 'vue';

const Tea = Vue.component('tea', {
    props: ['name', 'link'],
    template: `
      <li>
        <a :href="link" target="_blank">
          {{ name }}
        </a>
      </li>
    `
})

export default Tea;
