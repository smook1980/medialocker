<template lang="pug">
  div#app
    img(src="./assets/logo.png")
    h1 {{ msg }}
    el-button(@click.native="startHacking") Let's do it
</template>

<script>
export default {
  data() {
    return {
      msg: 'Use Vue 2.0 Today!',
    };
  },
  methods: {
    startHacking() {
      this.$notify(
        {
          title: 'It Works',
          message: 'We have laid the groundwork for you. Now it\'s your time to build something epic!',
          duration: 6000,
        }
      );
    },
  },
};
</script>

<style>
body {
  font-family: Helvetica, sans-serif;
}
</style>
