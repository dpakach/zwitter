export function getTimeStamp(dateString: String) : number {
  if (!dateString) {
    return 0
  }
  const date = dateString.split("-");
  var newDate = new Date( parseInt(date[0]), parseInt(date[1]) - 1, parseInt(date[2] ))
  return newDate.getTime()/1000
}

export function getDate(timeStamp: number): String {
  if (!timeStamp) {
    return ""
  }
  var date = new Date(timeStamp * 1000);
  var currentDate = date.toISOString().slice(0,10);
  return currentDate
}
