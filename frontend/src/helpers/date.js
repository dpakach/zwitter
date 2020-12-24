export function getTimeStamp(dateString) {
  if (!dateString) {
    return 0
  }
  const date = dateString.split("-");
  var newDate = new Date( date[0], date[1] - 1, date[2] );
  return newDate.getTime()/1000
}

export function getDate(timeStamp) {
  if (!timeStamp) {
    return 0
  }
  var date = new Date(timeStamp * 1000);
  var currentDate = date.toISOString().slice(0,10);
  return currentDate
}
