import React, { createRef, useState } from 'react'
import '../../style/crop.css'
import Cropper, { ReactCropperElement } from 'react-cropper'
import 'cropperjs/dist/cropper.css'
import { asAxiosError, generateURL, tryRequest } from '../../util'
import { useNavigate } from 'react-router-dom'
import axios from 'axios'
import config from '../../config'

export interface CropProperties {
	onClose: () => void
}

const ALLOWED_FILE_TYPES = ['image/png', 'image/jpeg', 'image/jpg']
const MIN_WIDTH = 150
const MIN_HEIGHT = 150
const MAX_FILE = 4_000_000
const MIN_FILE = 4_000

const Crop = function ({ onClose }: CropProperties) {
	const [imageSrc, setImageSrc] = useState('#')
	const [errorMessage, setErrorMessage] = useState('')
	const [uploadDisabled, setUploadDisabled] = useState(false)
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
			const src = reader.result as string
			const image = new Image()
			image.src = src
			image.onload = () => {
				if (
					image.naturalHeight < MIN_WIDTH ||
					image.naturalWidth < MIN_HEIGHT
				) {
					setImageSrc('#')
					setErrorMessage('Image too small')
				} else {
					setImageSrc(src)
					setErrorMessage('')
				}
			}
		}
		if (files !== null) reader.readAsDataURL(files[0])
	}

	const handleUpload = () => {
		if (imageSrc === '#') {
			return
		}
		setUploadDisabled(true)
		cropperRef.current?.cropper.getCroppedCanvas().toBlob((blob) => {
			const croppedData = blob
			tryRequest(
				async () => {
					console.log(croppedData)
					const formData = new FormData()
					if ((croppedData?.size as number) > MAX_FILE) {
						setErrorMessage('Cropped image too large')
						return
					}
					if ((croppedData?.size as number) < MIN_FILE) {
						setErrorMessage('Cropped image too small')
						return
					}
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
					setImageSrc('#')
					setUploadDisabled(false)
					location.reload()
				},
				(error) => {
					asAxiosError(error, (error) => {
						if (error.response?.status === 503) {
							setErrorMessage('Upload timeout')
						} else {
							console.log(error)
						}
					})
					setUploadDisabled(false)
				},
				navigate,
			)
		})
	}

	return (
		<div id='crop-page'>
			<div id='crop-topbar'>
				<label id='crop-input-button' htmlFor='crop-input'>
					Browse
				</label>
				<input
					id='crop-input'
					type='file'
					accept='image/png,image/jpg,image/jpeg'
					onChange={handleInputFileChange}
				/>
				<button id='crop-close-button' onClick={onClose}>
					Close
				</button>
			</div>
			<div className='image-area'>
				<div className='cropper'>
					<Cropper
						ref={cropperRef}
						style={{ aspectRatio: '1/1', height: '100%' }}
						zoomTo={0.1}
						aspectRatio={1}
						preview='.img-preview'
						background={false}
						src={imageSrc}
					/>
				</div>
				<div className='preview'>
					<div
						className='img-preview'
						style={{
							width: '100%',
							borderRadius: '50%',
							aspectRatio: '1/1',
							margin: '0',
						}}
					/>
				</div>
			</div>
			<div id='crop-error'>{errorMessage}</div>
			<button
				id='crop-upload-button'
				onClick={handleUpload}
				disabled={uploadDisabled}
			>
				Upload
			</button>
		</div>
	)
}

export default Crop
