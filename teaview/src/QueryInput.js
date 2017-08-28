import Vue from 'vue';
import { findTopTenTeaMatches } from './queryUtil';

const QueryInput = Vue.component('query-input', {
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
            this.$emit("update:newTeaList", matches.data)
          })
          .catch(error => console.error(error));
      }
    },
  },
  template: `
    <div>
      <input type="text" v-on:keyup="onKeyPress" />
      <button v-on:click="onSubmit">
        Submit
      </button>
    </div>
  `
})

export default QueryInput;
