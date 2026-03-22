package container

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TagSelection", func() {
	tags := []string{"v1.0.0", "v1.0.1", "v1.1.0", "v2.0.0-beta", "latest", "old-tag"}

	Describe("SelectBestTag", func() {
		Context("with semver strategy", func() {
			It("should pick the highest semver tag", func() {
				best, _ := SelectBestTag(tags, "semver", "")
				Expect(best).To(Equal("v2.0.0-beta"))
			})

			It("should pick the highest matching regex", func() {
				best, _ := SelectBestTag(tags, "semver", `^v1\.0\.\d+$`)
				Expect(best).To(Equal("v1.0.1"))
			})

			It("should handle human error (v1.1.0 pushed before v1.0.2)", func() {
				messyTags := []string{"v1.1.0", "v1.0.2"}
				best, _ := SelectBestTag(messyTags, "semver", "")
				Expect(best).To(Equal("v1.1.0"))
			})
		})

		Context("with invalid regex", func() {
			It("should return an error", func() {
				_, err := SelectBestTag(tags, "semver", "[")
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
