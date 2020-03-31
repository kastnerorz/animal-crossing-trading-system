export default function({ $axios, redirect }) {
  $axios.onRequest(config => {
    console.log('onRequest', config);
  })
  $axios.onError(error => {
    console.log('onError', error);
  });
}
