// currently not in use
<template>
  <div class="ac-btn">
    <div class="btn-icon"></div>
    <label></label>
  </div>
</template>

<script>
export default {
  props: {
    label: {
      type: String,
      required: true,
    },
    icon: {
      type: String,
      required: true,
    },
  },
}
</script>
