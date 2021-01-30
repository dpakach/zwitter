export function getTimeStamp(dateString: string): number {
  if (!dateString) {
    return 0;
  }
  const date = dateString.split('-');
  const newDate = new Date(parseInt(date[0]), parseInt(date[1]) - 1, parseInt(date[2]));
  return newDate.getTime() / 1000;
}

export function getDate(timeStamp: number): string {
  if (!timeStamp) {
    return '';
  }
  const date = new Date(timeStamp * 1000);
  const currentDate = date.toISOString().slice(0, 10);
  return currentDate;
}

export function timeSince(date: number): string {

  var seconds = Math.floor((new Date()).getTime() / 1000) - date;

  var interval = seconds / 31536000;

  if (interval > 1) {
    return Math.floor(interval) + " years";
  }
  interval = seconds / 2592000;
  if (interval > 1) {
    return Math.floor(interval) + " months";
  }
  interval = seconds / 86400;
  if (interval > 1) {
    return Math.floor(interval) + " days";
  }
  interval = seconds / 3600;
  if (interval > 1) {
    return Math.floor(interval) + " hours";
  }
  interval = seconds / 60;
  if (interval > 1) {
    return Math.floor(interval) + " minutes";
  }
  return Math.floor(seconds) + " seconds";
}