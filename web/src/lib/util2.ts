export function age2seconds(a: string) {
  if (!a && a.length < 2) return 0;
  let num = parseInt(a.slice(0, -1));
  switch (a.slice(-1)) {
    case 's':
      return num * 1;
    case 'm':
      return num * 60;
    case 'h':
      return num * 3600;
    case 'd':
      return num * 86400;
  }
  return 0;
}
