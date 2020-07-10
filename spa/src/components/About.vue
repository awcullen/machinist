<template>
  <b-jumbotron bg-variant="primary" text-variant="white" border-variant="dark">
    <template slot="header">
      Machinist Admin
    </template>
    <template slot="lead">
      Version: {{ version }}
    </template>
    <hr class="my-4">
    <p>
      <a href="https://golang.org/" target="_blank"><img src="../assets/go.svg" style="height:80px"/></a>&nbsp;
      <a href="https://vuejs.org/" target="_blank"><img src="../assets/logo.png" style="height:80px"/></a> 
    </p>
  </b-jumbotron>
</template>

<script>
const packageJson = require ('../../package.json')

export default {
  data: function() {
    return {
      version: packageJson.version
    }
  }
}
</script>
