import Vue from 'vue';
import axios from 'axios';

/**
 * need input for user to query
 * need way to display results
 * need way for user to click on results and see tea details that they got back
 */

function findTopTenTeaMatches(queryString) {
  return axios('http://127.0.0.1:5000/match', {
    method: 'POST',
    data: {
      userQuery: queryString,
    },
  });
}


Vue.component('tea', {
    props: ['name', 'link'],
    template: `
      <li>
        <a :href="link" target="_blank">
          {{ name }}
        </a>
      </li>
    `
})

Vue.component('tea-list', {
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
});

Vue.component('query-input', {
  data: function() {
    return {
      query: "",
    };
  },
  methods: {
    onKeyPress: function(event) {
      this.query = event.target.value;
    },
    onSubmit: function() {
      if (this.query) {
        findTopTenTeaMatches(this.query)
          .then(matches => {
            this.$emit("updateTeas", matches.data)
          })
          .catch(error => console.error(error));
      }
    },
  },
  template: `
    <div>
      <input type="text" v-on:keypress="onKeyPress" />
      <button v-on:click="onSubmit">
        Submit
      </button>
    </div>
  `
})

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
      <query-input v-on:updateTeas="updateTeaList" />
      <tea-list :teaList="teaList"></tea-list>
    </div>
  `
})
