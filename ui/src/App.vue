<template>
  <div id="app" class="uk-container">
    <div class="ml-nav-button">
      <a href="" class="uk-icon-button" uk-icon="icon: menu" uk-toggle="target: #offcanvas-overlay" ></a>
    </div>

    <!-- <button class="uk-button uk-button-default ml-nav-button" type="button" uk-toggle="target: #offcanvas-overlay">Open</button> -->
    <main>
      <img src="./assets/logo.png">
      <router-view></router-view>
    </main>

    <div id="offcanvas-overlay" uk-offcanvas="overlay: true">
      <div class="uk-offcanvas-bar">

      <h3>Title</h3>

      <p>Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</p>

      <button class="uk-button uk-button-default uk-offcanvas-close uk-width-1-1 uk-margin" type="button">Close</button>

      </div>
      </div>
  </div>
</template>

<script>
 export default {
   name: 'app',
 };
</script>

<style lang="less">
@import "./main.less";

#app {
  color: #2c3e50;
  margin-top: 60px;
}

#app main {
  /* font-family: 'Avenir', Helvetica, Arial, sans-serif; */
  text-align: center;
}

.ml-nav-button {
  position: absolute;
  text-align: left;
  width: 100%;
  top: 0;
  left: 0;
  padding-top: 1.5em;
  padding-left: 1.5em;
  margin: 0;
}

</style>
