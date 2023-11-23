
export function dateFormat(value) {
  if (value === '' || value === undefined) {
    return '';
  }
  const t = new Date(value);
  return t.getFullYear() + "-" + (t.getMonth() + 1) + "-" + t.getDate() + " " + t.getHours() + ":" + t.getMinutes() + ":" + t.getSeconds();
}
