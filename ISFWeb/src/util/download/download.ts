export const download = (content: string, fileName: string) => {
  // 加BOM(Byte Order Mark)使得csv文件使用UTF-8编码解析
  const BOM = '\uFEFF'
  const csvBlob = new Blob([BOM + content])
  if (navigator.msSaveBlob) { // IE10+
      navigator.msSaveBlob(csvBlob, fileName)
  } else {
      const csvURL = URL.createObjectURL(csvBlob)

      const anchor = document.createElement('a')
      anchor.href = csvURL
      anchor.download = fileName

      document.body.appendChild(anchor)
      anchor.click()
      document.body.removeChild(anchor)

      URL.revokeObjectURL(csvURL)
  }
}

/**
 * 从Content-Disposition头中提取文件名
 * @param headers 
 */
export const getFileName = (headers): string => {
  const contentDisposition = headers.get('content-disposition')
  
  let filename = 'download.xml'
  if (contentDisposition && contentDisposition.includes('filename=')) {
    const match = contentDisposition.match(/filename=([^;]+)/)
    if (match.length > 1) {
      filename = match[1]
    }
  }

  return filename
}