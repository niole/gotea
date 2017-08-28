import Vue from 'vue';
import Tea from './Tea';

const TeaList = Vue.component('tea-list', {
    props: ['teaList'],
    template: `
      <ul>
        <tea
          v-for="tea in teaList"
          :name="tea.name"
          :link="tea.link"
          :key="tea.link"
        />
      </ul>
    `
})

export default TeaList;
