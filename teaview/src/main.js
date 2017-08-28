import Vue from 'vue';
import QueryInput from './QueryInput';
import TeaList from './TeaList';

/**
 * Root of view tree
 */
new Vue({
  el: '#app',
  data: {
    teaList: [],
  },
  methods: {
    updateTeaList: function(newList) {
      this.teaList = newList;
    },
  },
  template: `
    <div>
      <query-input :newTeaList.sync="teaList" />
      <tea-list :teaList="teaList"></tea-list>
    </div>
  `
})
