const KiB = 1 << 10

export class Slice {
  constructor(private buf: Buffer = Buffer.allocUnsafe(4 * KiB), private len = 0) {
  }

  appendInPlace(elems: Uint8Array): Slice {
    const free = this.buf.length - this.len;
    // Grow buffer exponentially if elems won't fit into current buffer
    if (free < elems.length) {
      const oldBuf = this.buf
      this.buf = Buffer.allocUnsafe((this.len + elems.length) * 2);
      oldBuf.copy(this.buf)
    }
    Buffer.from(elems).copy(this.buf, this.len);
    this.len += elems.length
    return this
  }

  append(elems: Uint8Array): Slice {
    const slice = new Slice(this.buf, this.len)
    return slice.appendInPlace(elems)
  }

  buffer(): Buffer {
    return this.buf.slice(0, this.len)
  }

  length(): number {
    return this.len
  }

  slice(start = 0, end = this.len): Slice {
    if (end < 0) {
      end += this.len
    }
    return new Slice(this.buf.slice(start), end)
  }
}
