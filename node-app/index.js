import size from 'image-size'
import express from 'express'
import * as jose from 'jose'

const PORT = 8080
const HOST = '0.0.0.0'

const jwtSecret = new TextEncoder().encode('super-secret-key')
let requests = 0

size.disableFS(true)
const app = express()

app.use((req, res, next) => {
  requests++
  next()
})

app.use((req, res, next) => {
  const auth = req.headers['authorization']
  if (auth?.startsWith('Bearer ')) {
    const token = auth.substring(7)
    const decoded = jose
      .jwtVerify(token, jwtSecret)
      .then(decoded => {
        res.locals.user = decoded.payload.sub
      })
      .catch(next)
  }
  next()
})
app.get('/', (_, res) => {
  return res.send('OK')
})

app.use('/identify', express.raw({ limit: 1024000, type: '*/*' }))
app.post('/identify', (req, res) => {
  const buffer = req.body
  const sizeof = size(buffer)
  res.json({ ...sizeof, user: res.locals.user })
})

const generateJwt = (user = 'testuser') => {
  return new jose.SignJWT({})
    .setSubject(user)
    .setProtectedHeader({ alg: 'HS256' })
    .sign(jwtSecret)
}
app.listen(PORT, HOST, () => {
  generateJwt().then(jwt =>
    console.log(`Server is running on port ${PORT}, use JWT: ${jwt}`)
  )

  let prev = requests
  setInterval(() => {
    if (prev !== requests) {
      console.log('Requests: ' + requests)
      prev = requests
    }
  }, 1000)
})
