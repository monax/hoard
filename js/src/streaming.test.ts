import {Readable} from "stream";
import {pipeline} from "./pipeline";
import {bytesReadable, bytesToObject, objectToBytes, readAll, readBytes} from "./streaming";

type Buf = {
  buf: Uint8Array
}

const MiB = 1 << 20

describe('streaming', () => {
  test('buffering', async () => {
    const chunkSize = 100;
    const size = 50 * MiB
    const input = makeBuffer(size);
    const inputReadable = bytesReadable(input)
    const b2o = bytesToObject(buf => ({buf}), chunkSize)
    const objectPipeline = pipeline(inputReadable, b2o);
    const bufs: Buf[] = await readAll(objectPipeline)
    expect(bufs.length).toEqual(Math.ceil(size / chunkSize))
    const outputReadable = Readable.from(bufs)
    const o2b = objectToBytes<{ buf: Buffer }>(({buf}) => buf, chunkSize)
    const bytesPipeline = pipeline(outputReadable, o2b)
    const output = await readBytes(bytesPipeline)
    expect(output.equals(input)).toBeTruthy()
    expect(objectPipeline.destroyed).toBeTruthy()
    expect(bytesPipeline.destroyed).toBeTruthy()
  })
})

function makeBuffer(size: number): Buffer {
  const buf = Buffer.alloc(size)
  for (let i = 0; i < size; i++) {
    buf[i] = i % 256
  }
  return buf
}
