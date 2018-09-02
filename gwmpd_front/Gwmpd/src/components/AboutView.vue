<template lang="html">
  <div class="text-center topSpace">
    <b-container fluid class="mx-auto">
      <h2>Gwmpd</h2>
      <span class="text-primary">Freyermuth Julien</span><br>
      Copyright © 2018<br>
      <!-- <a href="https://en.wikipedia.org/wiki/BSD_licenses" target="_blank">BSD License</a> -->

      <a href="https://en.wikipedia.org/wiki/BSD_licenses" target="_blank"><b-img class="licenseImg" src="../assets/license.png" alt="Open Source logo"></b-img></a>

      <div class="donate">
        If you like my work, you can donate to support it. I will be even more motivated to add new features and improve the documentation<br>
        <form action="https://www.paypal.com/cgi-bin/webscr" method="post" target="_top" class="topSpace">
          <input type="hidden" name="cmd" value="_s-xclick">
          <input type="hidden" name="hosted_button_id" value="8SJ9D6CN28PSE">
          <input type="image" src="https://www.paypalobjects.com/fr_FR/FR/i/btn/btn_donateCC_LG.gif" border="0" name="submit" alt="PayPal, le réflexe sécurité pour payer en ligne">
          <img alt="" border="0" src="https://www.paypalobjects.com/fr_FR/i/scr/pixel.gif" width="1" height="1">
        </form>
      </div>

    </b-container>
  </div>
</template>

<script>
export default {
  name: 'about'
}
</script>

<style lang="scss" scoped>
  .donate{
    margin-top: 20px;
  }
  .licenseImg{
    width: 12%;
  }
  .topSpace{
    margin-top: 10px;
  }
</style>
