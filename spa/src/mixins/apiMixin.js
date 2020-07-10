export default {
  data: function () {
    return {
      apiEndpoint: "/api"
    }
  },
  
  methods: {
    apiGetHome: function() {  
      return fetch(`${this.apiEndpoint}/home`)
        .then(resp => {
          return resp.json();
        })
        .catch(err => {
          console.log(`### API Error! ${err}`);
        })
    },

    apiGetMetrics: function() {  
      return fetch(`${this.apiEndpoint}/metrics`)
        .then(resp => {
          return resp.json();
        })
        .catch(err => {
          console.log(`### API Error! ${err}`);
        })
    },

    apiGetInfo: function() {  
      return fetch(`${this.apiEndpoint}/info`)
        .then(resp => {
          return resp.json();
        })
        .catch(err => {
          console.log(`### API Error! ${err}`);
        })
    }    
  }
}