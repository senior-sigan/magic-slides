export function clamp(v: number, a:number, b?: number|undefined) {
  if (v < a) return a;
  if (b && v > b) return b;
  return v;
}