;;;; Copyright © 2022 siddharth ravikumar <s@ricketyspace.net>
;;;; SPDX-License-Identifier: ISC

#lang racket

(require racket/draw)
(require racket/snip)

;; Sets up settings for the drawing context.
(define (setup-dc dc)
  (send dc set-smoothing 'aligned)
  (send dc set-pen "black" 1 'transparent)
  (send dc set-brush "black" 'solid))

;; Draws a moon.
(define (moon-drawing size)
  (let* ((moon (new dc-path%))
         (inner-arc-x (/ (* size -50) 500.0))
         (inner-arc-y (* -1 inner-arc-x))
         (inner-arc-size (/ (* size 320) 500.0))
         (outer-arc-x (/ (* size 5) 500.0))
         (outer-arc-y outer-arc-x)
         (outer-arc-size (/ (* size 490) 500.0)))
    (send moon arc inner-arc-x inner-arc-y inner-arc-size inner-arc-size 1.57 4.36 #f)
    (send moon arc outer-arc-x outer-arc-y outer-arc-size outer-arc-size 3.54 2.20 #t)
    moon))

;; Draws the peach logo in a bitmap and returns the bitmap.
(define (draw-logo size)
  (let* ((target (make-bitmap size size))
         (dc (new bitmap-dc% (bitmap target))))
    (setup-dc dc)
    (send dc set-brush "black" 'solid)
    (send dc draw-path (moon-drawing size))
    target))

;; Exports the logo into PNG.
(define (png-export logo size)
  (send logo save-file (format "peach-~s.png" size) 'png))

;; Peach logo as a bitmap.
(define peach-logo (draw-logo 500))

;; Render logo in racket shell.
(make-object image-snip% peach-logo)

;; Export logo to a PNG file.
(png-export peach-logo 500)
