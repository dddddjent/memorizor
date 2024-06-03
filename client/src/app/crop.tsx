import React, { createRef, useState } from 'react'
import '../style/crop.css'
import Cropper, { ReactCropperElement } from 'react-cropper'
import 'cropperjs/dist/cropper.css'
import { generateURL, tryRequest } from '../util'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import config from '../config'

export interface CropProperties {
    onClose: () => void
}

const ALLOWED_FILE_TYPES = ['image/png', 'image/jpeg', 'image/jpg']

const Crop = function({ onClose }: CropProperties) {
    const [imageSrc, setImageSrc] = useState('#')
    const cropperRef = createRef<ReactCropperElement>()
    const navigate = useNavigate()

    const handleInputFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        e.preventDefault()
        const files = e.target.files
        if (files == null) return

        const allAllowed = Array.from(files).every((file) =>
            ALLOWED_FILE_TYPES.includes(file.type),
        )
        if (!allAllowed) return

        const reader = new FileReader()
        reader.onload = () => {
            setImageSrc(reader.result as string)
        }
        if (files !== null) reader.readAsDataURL(files[0])
    }

    const handleUpload = () => {
        if (imageSrc === '#') {
            return
        }
        let croppedData: Blob | null
        cropperRef.current?.cropper.getCroppedCanvas().toBlob((blob) => {
            croppedData = blob
        })
        tryRequest(
            async () => {
                const formData = new FormData()
                formData.append('image', croppedData as Blob)
                const { data } = await axios.post(
                    generateURL(config.api.account, '/profile_image'),
                    formData,
                    {
                        headers: {
                            'Content-Type': 'multipart/form-data',
                        },
                    },
                )
                console.log(data)
                location.reload()
                setImageSrc('#')
            },
            (error) => {
                console.log(error)
            },
            navigate,
        )
    }

    return (
        <div id='crop-page'>
            <div>
                <input
                    type='file'
                    accept='image/png,image/jpg,image/jpeg'
                    onChange={handleInputFileChange}
                />
                <button onClick={onClose}>Close</button>
            </div>
            <div>
                <Cropper
                    ref={cropperRef}
                    style={{ width: '20%' }}
                    zoomTo={0.1}
                    initialAspectRatio={1}
                    preview='.img-preview'
                    src={imageSrc}
                />
            </div>
            <div>
                <button onClick={handleUpload}>Upload</button>
            </div>
        </div>
    )
}

export default Crop
