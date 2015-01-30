package obj_test

import (
	"io"
	"os"

	. "github.com/momchil-atanasov/go-data-front/obj"
	"github.com/momchil-atanasov/go-data-front/obj/test_helpers"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Scanner", func() {

	var handlerFixture *test_helpers.HandlerFixture
	var scanner Scanner
	var scanError error

	BeforeEach(func() {
		handlerFixture = test_helpers.NewHandlerFixture()
		scanner = NewScanner(handlerFixture.Handler())
		scanError = nil
	})

	scan := func(reader io.Reader) {
		scanError = scanner.Scan(reader)
	}

	scanFile := func(filename string) {
		file, err := os.Open(filename)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		scan(file)
	}

	itShouldHaveReturnedAnError := func() {
		It("should have returned an error", func() {
			Ω(scanError).Should(HaveOccurred())
		})
	}

	itShouldNotHaveReturnedAnError := func() {
		It("should not have returned an error", func() {
			Ω(scanError).ShouldNot(HaveOccurred())
		})
	}

	Context("when a basic OBJ file is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_basic.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the comments", func() {
			handlerFixture.AssertCommentCall("This is the beginning of this OBJ file.")
			handlerFixture.AssertNoMoreCommentCalls()
		})

		It("should have scanned the vertices", func() {
			handlerFixture.AssertVertexXYZ(-1.0, 1.0, -1.0)
			handlerFixture.AssertVertexXYZ(-1.0, -1.0, 1.0)
			handlerFixture.AssertVertexXYZ(1.0, -1.0, -1.0)
			handlerFixture.AssertVertexXYZ(1.0, 1.0, 1.0)
			handlerFixture.AssertNoMoreVertices()
		})

		It("should have scanned the texture coordinates", func() {
			handlerFixture.AssertTexCoordUV(0.0, 0.0)
			handlerFixture.AssertTexCoordUV(1.0, 1.0)
			handlerFixture.AssertTexCoordUV(1.0, 0.0)
			handlerFixture.AssertTexCoordUV(0.0, 1.0)
		})

		It("should have scanned the normals", func() {
			handlerFixture.AssertNormal(0.0, 1.0, 0.0)
			handlerFixture.AssertNormal(1.0, 0.0, 0.0)
			handlerFixture.AssertNormal(0.0, 0.0, 1.0)
		})

		It("should have scanned the objects", func() {
			handlerFixture.AssertObject("MyObject")
		})

		It("should have scanned the faces", func() {
			handlerFixture.AssertFaceStart()
			handlerFixture.AssertFaceEnd()
		})

		It("should have scanned data references", func() {
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(1)
			handlerFixture.AssertTexCoordIndex(4)
			handlerFixture.AssertNormalIndex(1)
			handlerFixture.AssertCoordReferenceEnd()

			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(2)
			handlerFixture.AssertTexCoordIndex(1)
			handlerFixture.AssertNormalIndex(1)
			handlerFixture.AssertCoordReferenceEnd()

			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(3)
			handlerFixture.AssertTexCoordIndex(3)
			handlerFixture.AssertNormalIndex(2)
			handlerFixture.AssertCoordReferenceEnd()
		})

		It("should have scanned material libraries", func() {
			handlerFixture.AssertMaterialLibrary("valid_basic.mtl")
		})

		It("should have scanned material references", func() {
			handlerFixture.AssertMaterialReference("BlueMaterial")
		})
	})

	Context("when a file with all kinds of comments is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_comments.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the comments", func() {
			handlerFixture.AssertCommentCall("Comment at file start")
			handlerFixture.AssertCommentCall("Comment that is right next to special char")
			handlerFixture.AssertCommentCall("This comment uses")
			handlerFixture.AssertCommentCall("two lines")
			handlerFixture.AssertCommentCall("")
			handlerFixture.AssertCommentCall("Previous comment was empty. This one contain the # character twice.")
			handlerFixture.AssertCommentCall("Comment at file end")
			handlerFixture.AssertNoMoreCommentCalls()
		})
	})

	Context("when a file with all kinds of vertices is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_vertices.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the vertices", func() {
			handlerFixture.AssertVertexXYZ(1.0, 1.0, -1.0)
			handlerFixture.AssertVertexXYZW(-1.0, -1.0, 1.0, 0.5)
			handlerFixture.AssertNoMoreVertices()
		})
	})

	Context("when a file with all kinds of texture coordinates is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_texcoords.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the texture coordinates", func() {
			handlerFixture.AssertTexCoordU(1.6)
			handlerFixture.AssertTexCoordUV(0.0, -0.5)
			handlerFixture.AssertTexCoordUVW(-0.2, 1.4, 3.0)
			handlerFixture.AssertNoMoreTexCoords()
		})
	})

	Context("when a file with all kinds of normals is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_normals.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the normals", func() {
			handlerFixture.AssertNormal(1.0, 0.0, 0.0)
			handlerFixture.AssertNormal(-1.0, -1.0, 1.0)
			handlerFixture.AssertNoMoreNormals()
		})
	})

	Context("when a file with all kinds of objects is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_objects.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the objects", func() {
			handlerFixture.AssertObject("FirstObject")
			handlerFixture.AssertObject("LastObject")
			handlerFixture.AssertNoMoreObjects()
		})
	})

	Context("when a file with all kinds of faces is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_faces.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned them correctly", func() {
			handlerFixture.AssertFaceCallCount(3)
			handlerFixture.AssertCoordReferenceCallCount(12)
		})
	})

	Context("when a file with all kinds of coord references is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_coord_references.obj")
		})

		It("should not have returned an error", func() {
			Ω(scanError).ShouldNot(HaveOccurred())
		})

		It("should have scanned them correctly", func() {
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(1)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(2)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(3)
			handlerFixture.AssertCoordReferenceEnd()

			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(2)
			handlerFixture.AssertTexCoordIndex(1)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(3)
			handlerFixture.AssertTexCoordIndex(2)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(4)
			handlerFixture.AssertTexCoordIndex(3)
			handlerFixture.AssertCoordReferenceEnd()

			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(4)
			handlerFixture.AssertNormalIndex(5)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(6)
			handlerFixture.AssertNormalIndex(7)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(8)
			handlerFixture.AssertNormalIndex(9)
			handlerFixture.AssertCoordReferenceEnd()

			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(1)
			handlerFixture.AssertTexCoordIndex(2)
			handlerFixture.AssertNormalIndex(3)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(2)
			handlerFixture.AssertTexCoordIndex(3)
			handlerFixture.AssertNormalIndex(4)
			handlerFixture.AssertCoordReferenceEnd()
			handlerFixture.AssertCoordReferenceStart()
			handlerFixture.AssertVertexIndex(3)
			handlerFixture.AssertTexCoordIndex(4)
			handlerFixture.AssertNormalIndex(5)
			handlerFixture.AssertCoordReferenceEnd()
		})
	})

	Context("when a file with all kinds of material libraries is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_material_libraries.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the material libraries", func() {
			handlerFixture.AssertMaterialLibrary("valid.mtl")
			handlerFixture.AssertMaterialLibrary("extension01.mtl")
			handlerFixture.AssertMaterialLibrary("extension02.mtl")
			handlerFixture.AssertNoMoreMaterialLibraries()
		})
	})

	Context("when a file with all kinds of material references is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_material_references.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned the material references", func() {
			handlerFixture.AssertMaterialReference("")
			handlerFixture.AssertMaterialReference("MyMaterial")
			handlerFixture.AssertNoMoreMaterialReferences()
		})
	})

	Context("when a file with complex logical line is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_logical_line.obj")
		})

		itShouldNotHaveReturnedAnError()

		It("should have scanned everything correctly", func() {
			handlerFixture.AssertVertexXYZ(-1.0, 1.0, -2.0)
			handlerFixture.AssertVertexXYZW(-1.0, -1.0, 1.0, 2.0)
			handlerFixture.AssertNoMoreVertices()

			handlerFixture.AssertObject("HelloWorld")
			handlerFixture.AssertObject("NonNewLine")
			handlerFixture.AssertObject("EndsWithNewLine")
			handlerFixture.AssertNoMoreObjects()
		})
	})

	Context("when a file with unknown command is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/valid_unknown_command.obj")
		})

		itShouldNotHaveReturnedAnError()
	})

	Context("when a file with insufficient vertex data is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/error_insufficient_vertex_data.obj")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with insufficient texture coordinate data is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/error_insufficient_texcoord_data.obj")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with insufficient normal data is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/error_insufficient_normal_data.obj")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with an unnamed object is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/error_empty_object_name.obj")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with corrupt vertex is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/error_corrupt_vertex.obj")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with corrupt texture coordinate is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/error_corrupt_texcoord.obj")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with corrupt normal is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/error_corrupt_normal.obj")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with corrupt vertex reference is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/error_corrupt_vertex_reference.obj")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with corrupt texture coordinate reference is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/error_corrupt_texcoord_reference.obj")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when a file with corrupt normal reference is scanned", func() {
		BeforeEach(func() {
			scanFile("testres/error_corrupt_normal_reference.obj")
		})

		itShouldHaveReturnedAnError()
	})

	Context("when reading fails", func() {
		BeforeEach(func() {
			reader := new(test_helpers.FailingReader)
			scan(reader)
		})

		itShouldHaveReturnedAnError()
	})
})
