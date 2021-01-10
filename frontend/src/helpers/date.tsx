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
