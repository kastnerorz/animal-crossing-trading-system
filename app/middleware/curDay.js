export default function({ store, route, redirect }) {
  const cDate = new Date()
  // cDate.setDate(cDate.getDate()+2);
  const cDay = cDate.getDay()
  const cHour = cDate.getHours()
  // if (route.name === 'index') {
  //   if (cDay !== 0 && cHour < 12) {
  //     redirect('/takeRest?type=buy')
  //   }
  // } else if (route.name === 'sell') {
  //   if (cDay === 0) {
  //     redirect('/takeRest?type=sell')
  //   }
  // }
}
