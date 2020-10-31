package envelopespec_test

import (
	"time"

	. "github.com/dogmatiq/marshalkit/fixtures"

	"github.com/dogmatiq/dogma"
	. "github.com/dogmatiq/dogma/fixtures"
	. "github.com/dogmatiq/envelopespec"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("func CheckWellFormed()", func() {
	var env *Envelope

	BeforeEach(func() {
		createdAt := time.Now()
		scheduledFor := createdAt.Add(1 * time.Hour)

		env = &Envelope{
			MessageId:     "<id>",
			CausationId:   "<cause>",
			CorrelationId: "<correlation>",
			SourceApplication: &Identity{
				Name: "<app-name>",
				Key:  "<app-key>",
			},
			SourceHandler: &Identity{
				Name: "<handler-name>",
				Key:  "<handler-key>",
			},
			SourceInstanceId: "<instance>",
			CreatedAt:        MarshalTime(createdAt),
			ScheduledFor:     MarshalTime(scheduledFor),
			Description:      dogma.DescribeMessage(MessageA1),
			PortableName:     MessageAPortableName,
			MediaType:        MessageA1Packet.MediaType,
			Data:             MessageA1Packet.Data,
		}
	})

	It("does not return an error if the envelope is well-formed", func() {
		err := CheckWellFormed(env)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("returns an error if the message ID is empty", func() {
		env.MessageId = ""

		err := CheckWellFormed(env)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the causation ID is empty", func() {
		env.CausationId = ""

		err := CheckWellFormed(env)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the correlation ID is empty", func() {
		env.CorrelationId = ""

		err := CheckWellFormed(env)
		Expect(err).Should(HaveOccurred())
	})

	It("returns an error if the source app name is empty", func() {
		env.SourceApplication.Name = ""

		err := CheckWellFormed(env)
		Expect(err).To(MatchError("application identity is invalid: identity name must not be empty"))
	})

	It("returns an error if the source app key is empty", func() {
		env.SourceApplication.Key = ""

		err := CheckWellFormed(env)
		Expect(err).To(MatchError("application identity is invalid: identity key must not be empty"))
	})

	It("returns an error if the source handler name is empty", func() {
		env.SourceHandler.Name = ""

		err := CheckWellFormed(env)
		Expect(err).To(MatchError("handler identity is invalid: identity name must not be empty"))
	})

	It("returns an error if the source handler key is empty", func() {
		env.SourceHandler.Key = ""

		err := CheckWellFormed(env)
		Expect(err).To(MatchError("handler identity is invalid: identity key must not be empty"))
	})

	It("returns an error if the source handler is empty but the instance ID is set", func() {
		env.SourceHandler = nil

		err := CheckWellFormed(env)
		Expect(err).To(MatchError("source instance ID must not be specified without providing a handler identity"))
	})

	When("there is no source handler", func() {
		BeforeEach(func() {
			env.SourceHandler = nil
			env.SourceInstanceId = ""
		})

		It("returns an error if the message is a timeout", func() {
			err := CheckWellFormed(env)
			Expect(err).To(MatchError("scheduled-for time must not be specified without a providing source handler and instance ID"))
		})

		It("does not return an error if the message is not a timeout", func() {
			env.ScheduledFor = ""

			err := CheckWellFormed(env)
			Expect(err).ShouldNot(HaveOccurred())
		})
	})

	It("does not return an error if the message description is empty", func() {
		env.Description = ""

		err := CheckWellFormed(env)
		Expect(err).ShouldNot(HaveOccurred())
	})

	It("returns an error if the created-at timestamp is empty", func() {
		env.CreatedAt = ""

		err := CheckWellFormed(env)
		Expect(err).To(MatchError("created-at time must not be empty"))
	})

	It("returns an error if the portable name is empty", func() {
		env.PortableName = ""

		err := CheckWellFormed(env)
		Expect(err).To(MatchError("portable name must not be empty"))
	})

	It("returns an error if the media-type is empty", func() {
		env.MediaType = ""

		err := CheckWellFormed(env)
		Expect(err).To(MatchError("media-type must not be empty"))
	})

	It("does not return an error if the message data is empty", func() {
		env.Data = nil

		err := CheckWellFormed(env)
		Expect(err).ShouldNot(HaveOccurred())
	})
})

var _ = Describe("func MustBeWellFormed()", func() {
	It("panics if the envelope is not well-formed", func() {
		Expect(func() {
			MustBeWellFormed(&Envelope{})
		}).To(Panic())
	})
})