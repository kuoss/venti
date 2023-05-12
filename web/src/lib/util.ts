import ColorHash from 'color-hash';
import yaml from 'js-yaml';

const colorHash = new ColorHash({
  lightness: [0.4, 0.55, 0.7],
  saturation: [0.4, 0.55, 0.7],
});

export default {
  utc2local(s: string) {
    const d = new Date(s)
    return d.toLocaleDateString('fr-CA') + ' ' + d.toTimeString().slice(0, 8);
  },
  dateTimeAsLocal(d: number) {
    let dt = new Date(d * 1000);
    dt = new Date(dt.getTime() - dt.getTimezoneOffset() * 60 * 1000);
    return dt.toISOString().substring(0, 19).replace('T', ' ');
  },
  utc2age(u: number) {
    return this.seconds2age((new Date().getTime() - new Date(u).getTime()) / 1000);
  },
  age2seconds(a: string) {
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
  },
  seconds2age(s: number) {
    if (s == 0) return 0;
    if (s > 86400) return Math.floor(s / 86400) + 'd';
    if (s > 3600) return Math.floor(s / 3600) + 'h';
    if (s > 60) return Math.floor(s / 60) + 'm';
    return Math.floor(s) + 's';
  },
  nanoseconds2human(ns: string) {
    const num = parseInt(ns);
    return this.seconds2age(num / 1000000000);
  },
  string2color(str: string) {
    return colorHash.hex(str);
  },
  string2letters(str: string) {
    return str.split(/[.-]+/).map((x) => x[0]).join('');
  },
  copyToClipboard(val: string) {
    const t = document.createElement('textarea');
    document.body.appendChild(t);
    t.value = val;
    t.select();
    document.execCommand('copy'); // TODO: execCommand deprecated
    document.body.removeChild(t);
  },
  dumpYAML1(j: string, flowLevel: number) {
    return yaml
      .dump(j, { noArrayIndent: true, flowLevel: flowLevel })
      .replace(/>-/g, '|')
      .replace(/>/g, '|');
  },
  dumpYAML2(s: string) {
    return yaml.dump(s)
      .replace(/>-/g, '|')
      .replace(/>/g, '|')
      .replace(/\n/g, '\n  ')
      .trimEnd()
  },
  indentText(t: string, level: number): string {
    return t.split('\n').map((x) => ' '.repeat(level) + x).join('\n').trimEnd();
  },
  cloneObject(o: Object): Object {
    return JSON.parse(JSON.stringify(o));
  },

}