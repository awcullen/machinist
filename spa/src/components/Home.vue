<template>
  <div class="card">
    <div class="card-header bg-success h3">
      <fa icon="home" />&nbsp; Home
    </div>
    <div class="card-body">
      <spinner v-if="!info"></spinner>
      <div v-if="info">
        <b>Gateway</b>
        <table class="table">
          <colgroup>
            <col span="1" style="width: 25%;">
            <col span="1" style="width: 75%;">
          </colgroup>
          <tbody>
            <tr>
              <td><b>Name</b></td>
              <td>{{ info.name }}</td>
            </tr>
            <tr>
              <td><b>Location</b></td>
              <td>{{ info.location }}</td>
            </tr>
            <tr>
              <td><b>Project Id</b></td>
              <td>{{ info.projectId }}</td>
            </tr>
            <tr>
              <td><b>Region</b></td>
              <td>{{ info.region }}</td>
            </tr>
            <tr>
              <td><b>Registry Id</b></td>
              <td>{{ info.registryId }}</td>
            </tr>
            <tr>
              <td><b>Device Id</b></td>
              <td>{{ info.deviceId }}</td>
            </tr>
            <tr>
              <td><b>Store Path</b></td>
              <td>{{ info.storePath }}</td>
            </tr>
          </tbody>        
        </table>
        <b>Machines</b>
        <table v-for="device in info.devices" :key="device.deviceID" class="table">
          <colgroup>
            <col span="1" style="width: 25%;">
            <col span="1" style="width: 75%;">
          </colgroup>
           <tbody>
            <tr>
              <td><b>Name</b></td>
              <td>{{ device.name }}</td>
            </tr>
            <tr>
              <td><b>Location</b></td>
              <td>{{ device.location }}</td>
            </tr>
            <tr>
              <td><b>Device Id</b></td>
              <td>{{ device.deviceId }}</td>
            </tr>
            <tr>
              <td><b>Kind</b></td>
              <td>{{ device.kind }}</td>
            </tr>
          </tbody>        
        </table>
      </div>
    </div>
  </div>
</template>

<script>
import apiMixin from "../mixins/apiMixin.js";
import Spinner from "./Spinner.vue";
const info = null;

export default {
  mixins: [apiMixin],

  data: function() {
    return {
      info: info
    };
  },

  components: {
    Spinner
  },

  created() {
    this.getHome();
  },

  methods: {
    getHome: function() {
      fetch(`${this.apiEndpoint}/home`)
        .then(resp => {
          return resp.json();
        })
        .then(json => {
          this.info = json;
        })
        .catch(err => {
          console.log(err);
        })
    }
  },

  filters: {
    titleify: function(value) {
      if (!value) return "";
      value = value.toString();
      value = value.replace(/([A-Z])/g, ' $1')
      value = value.replace(/^./, function(str){ return str.toUpperCase(); });
      return value;
    }
  },

  computed: {
  }  
}
</script>
