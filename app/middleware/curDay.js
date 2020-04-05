export default function({ store, route, redirect }) {
  const cDate = new Date();
  // if (process.env.NODE_ENV === "development") {
  //   cDate.setDate(cDate.getDate() + 2);
  // }
  const cDay = cDate.getDay();
  const cHour = cDate.getHours();
  if (route.name === "sell") {
    if (cDay === 0) {
      redirect("/");
    }
  } else if (route.name === "index") {
    if (cDay !== 0) {
      redirect("/sell");
    } else if (cHour > 11) {
      redirect("/takeRest");
    }
  } else if (route.name === "application") {
    if (cDay === 0 && cHour > 11) {
      redirect("/takeRest");
    }
  } else if (route.name === "takeRest") {
    if (cDay !== 0) {
      redirect("/sell");
    } else if (cHour < 12) {
      redirect("/");
    }
  }
}
