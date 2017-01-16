<template lang="pug">
  div#app
    el-tabs(v-model='activeName', @tab-click='handleClick')
      el-tab-pane(label='User', name='first') User
      el-tab-pane(label='Config', name='second') Config
      el-tab-pane(label='Role', name='third') Role
      el-tab-pane(label='Task', name='fourth') Task
</template>

<script>
export default {
  data() {
    return {
      activeName: 'first',
    };
  },
  methods: {
    handleClick(tab, event) {
      console.log(tab, event);
    },
  },
};
</script>

<style>
 /* body {
    font-family: 'Open Sans', sans-serif;
    }
  */
</style>
